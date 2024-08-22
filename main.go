package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	//making sure the correct amount of inputs is given
	if len(os.Args) != 3 {
		fmt.Println("Incorrect amount of inputs :(, try again with the correct amount!")
		return
	}
	inputfilename := os.Args[1]  //taking the name of the file to read and storing it
	outputfilename := os.Args[2] //taking the second input of the file and storing it in outputfilename

	if os.Args[1] != "sample.txt" || os.Args[2] != "result.txt" {
		fmt.Println("Error: invalid file names, only sample.txt and results.txt allowed in their respective orders")
		os.Exit(0)
	}

	file, err := os.Open(inputfilename) //opening the file
	if err != nil {
		fmt.Println("There appears to be an error in opening the file, specifics are as follows : ", err)
		return
	}
	defer file.Close()
	//closes the file when its done with processes
	content, err := io.ReadAll(file)
	//stores the file in content [] byte
	if err != nil {
		fmt.Println("There appears to be an error reading the file, specifics are as follows: ", err)
		return
	}

	modifiedContent := modifyContent(string(content)) //invokes modifycontent which will call all functions needed and return the final answer
	fmt.Println("Modified Content:")
	fmt.Println(modifiedContent)

	// Write the modified content to the output file
	err = writeToFile(outputfilename, modifiedContent)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Successfully wrote modified content to", outputfilename)
}

func modifyContent(s string) string {
	// Define patterns for each type of transformation
	hexPattern := regexp.MustCompile(`([0-9A-Fa-f]+)\s*\(hex\)`)
	binPattern := regexp.MustCompile(`([01]+)\s*\(bin\)`)
	upPattern := regexp.MustCompile(`((\w+\s*)+)\s*\(up(?:,\s*(\d+))?\)`)
	lowPattern := regexp.MustCompile(`((\w+\s*)+)\s*\(low(?:,\s*(\d+))?\)`)
	capPattern := regexp.MustCompile(`((\w+\s*)+)\s*\(cap(?:,\s*(\d+))?\)`)

	// Apply hex conversion
	if hexPattern.MatchString(s) {
		s = hexConvert(s)
	}

	// Apply binary conversion
	if binPattern.MatchString(s) {
		s = binConvert(s)
	}

	// Apply uppercase conversion
	if upPattern.MatchString(s) {
		s = upConvert(s)
	}

	// Apply lowercase conversion
	if lowPattern.MatchString(s) {
		s = lowConvert(s)
	}

	// Apply capitalization
	if capPattern.MatchString(s) {
		s = capConvert(s)
	}

	// Format text
	s = formatText(s)

	return s
}
func hexConvert(text string) string {
	hexPattern := regexp.MustCompile(`([0-9A-Fa-f]+)\s*\(hex\)`)
	//checking the pattern of words followed by a (hex) identifier
	//make sure to use the spicy single quotations near the 1, saved me a lot of trouble
	matches := hexPattern.FindAllStringSubmatch(text, -1)
	//checking all instances of (hex)
	fmt.Println("Match results:", matches)
	//displaying our finds :)

	modifiedtext := text
	for _, match := range matches {
		hexvalue := match[1]
		//saving hex value
		decimalValue, err := strconv.ParseInt(hexvalue, 16, 64)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		modifiedtext = strings.ReplaceAll(modifiedtext, match[0], fmt.Sprintf("%d", decimalValue))

	}
	return modifiedtext
}

func binConvert(text string) string {
	binPattern := regexp.MustCompile(`([01]+)\s*\(bin\)`)
	//checking the pattern of words followed by a (hex) identifier
	//make sure to use the spicy single quotations near the 1, saved me a lot of trouble
	matches := binPattern.FindAllStringSubmatch(text, -1)
	//checking all instances of (hex)
	fmt.Println("Bin Match results:", matches)
	//displaying our finds :)

	modifiedtext := text
	for _, match := range matches {
		binvalue := match[1]
		//saving hex value
		decimalValue, err := strconv.ParseInt(binvalue, 2, 64)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		modifiedtext = strings.ReplaceAll(modifiedtext, match[0], fmt.Sprintf("%d", decimalValue))

	}
	return modifiedtext
}

func upConvert(text string) string {
	upPattern := regexp.MustCompile(`((\w+\s*)+)\s*\(up(?:,\s*(\d+))?\)`)
	// regular expression pattern meaning 1 or more words with optional spaces followed by (up) and optionally followed by , n; n being a number
	matches := upPattern.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		words := strings.Fields(match[1]) //splits the words, based on the first capturing group, which is the words split after whitespaces
		numWords := len(words)            //gets the number of words found in the match
		n := numWords                     //Initializes n to the total number of words, meaning all words will be transformed by default.
		if len(match) > 3 && match[3] != "" {
			n, _ = strconv.Atoi(match[3]) // Get the number after ", "
		}
		if n > numWords { //Ensures n doesn’t exceed the actual number of words in case a larger number was specified.
			n = numWords
		}

		for i := numWords - n; i < numWords; i++ {
			words[i] = strings.ToUpper(words[i]) //Loops through the first n words and converts each one to uppercase using strings.ToUpper().
		}
		//Uses strings.Join() to combine the transformed words back into a single string.
		//Replaces the original part of text (indicated by match[0]) with the updated string.

		text = strings.ReplaceAll(text, match[0], strings.Join(words, " "))
	}
	return text

}

func lowConvert(text string) string {
	lowPattern := regexp.MustCompile(`((\w+\s*)+)\s*\(low(?:,\s*(\d+))?\)`)
	// regular expression pattern meaning 1 or more words with optional spaces followed by (low) and optionally followed by , n; n being a number
	matches := lowPattern.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		words := strings.Fields(match[1]) //splits matched words into slice of strings
		numWords := len(words)            //gets the number of words found in the match
		n := numWords                     //Initializes n to the total number of words, meaning all words will be transformed by default.
		if len(match) > 3 && match[3] != "" {
			n, _ = strconv.Atoi(match[3]) // Get the number after ", "
		}
		if n > numWords { //Ensures n doesn’t exceed the actual number of words in case a larger number was specified.
			n = numWords
		}

		for i := 0; i < n; i++ {
			words[i] = strings.ToLower(words[i]) //Loops through the first n words and converts each one to uppercase using strings.ToUpper().
		}
		//Uses strings.Join() to combine the transformed words back into a single string.
		//Replaces the original part of text (indicated by match[0]) with the updated string.

		text = strings.ReplaceAll(text, match[0], strings.Join(words, " "))
	}
	return text

}

func capConvert(text string) string {
	capPattern := regexp.MustCompile(`((\w+\s*)+)\s*\(cap(?:,\s*(\d+))?\)`)

	matches := capPattern.FindAllStringSubmatch(text, -1) // find all submatches

	for _, match := range matches {
		words := strings.Fields(match[1])
		numWords := len(words)
		n := numWords
		// Check for a specific number
		if len(match) > 3 && match[3] != "" {
			n, _ = strconv.Atoi(match[3]) // Get the number after ", "
		}

		// Ensure n doesn't exceed the number of words
		if n > numWords {
			n = numWords
		}

		// Capitalize words
		for i := 0; i < n; i++ {
			words[i] = capitalize(words[i])
		}

		// Replace the matched part in text with the updated words
		text = strings.ReplaceAll(text, match[0], strings.Join(words, " "))
	}
	return text

}
func capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	return strings.ToUpper(string(word[0])) + word[1:]
}

func formatText(text string) string {
	// 1. Remove spaces before specified punctuation marks
	text = regexp.MustCompile(`\s+([.,!?;:])`).ReplaceAllString(text, "$1")

	// 2. Ensure there is a single space after commas
	text = regexp.MustCompile(`,\s*`).ReplaceAllString(text, ", ")

	// 3. Ensure no space after other punctuation marks
	text = regexp.MustCompile(`([!?;:])\s*`).ReplaceAllString(text, "$1")

	// 4. Handle ellipses
	text = regexp.MustCompile(`\s*\.\.\.\s*`).ReplaceAllString(text, "... ")

	// 5. Remove spaces around single quotes
	text = regexp.MustCompile(`\s*'\s*`).ReplaceAllString(text, "'")

	// 6. Support various punctuation groups
	text = regexp.MustCompile(`!{2,}|[!?]{2,}`).ReplaceAllString(text, "$0")

	// 7. Change "a" to "an" before words starting with vowels or 'h'
	text = TransformArticles(text)

	return text
}
func TransformArticles(text string) string {
	// Change "a" to "an" before words starting with a vowel or silent 'h'
	text = regexp.MustCompile(`\ba\s+([aeiouAEIOU]|h[aeiou])`).ReplaceAllString(text, "an $1")
	return text
}

func writeToFile(filename, content string) error {
	file, err := os.Create(filename) // Create the output file
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file is closed when done

	_, err = io.WriteString(file, content) // Write the content
	return err
}
