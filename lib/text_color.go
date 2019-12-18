// the library of text color, set the color to text tag of message.
package lib

import "fmt"

// number value constant of color.
const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

// Black color of black.
func Black(str string) string {
	return textColor(TextBlack, str)
}

// Red color of Red, usually use in the type of error, fatal message.
func Red(str string) string {
	return textColor(TextRed, str)
}

// Green color of Green.
func Green(str string) string {
	return textColor(TextGreen, str)
}

// Yellow color of Yellow.
func Yellow(str string) string {
	return textColor(TextYellow, str)
}

// Blue color of Blue.
func Blue(str string) string {
	return textColor(TextBlue, str)
}

// Magenta color of Magenta.
func Magenta(str string) string {
	return textColor(TextMagenta, str)
}

// Cyan color of Cyan.
func Cyan(str string) string {
	return textColor(TextCyan, str)
}

// White color of White.
func White(str string) string {
	return textColor(TextWhite, str)
}

// textColor set color of message tag.
func textColor(color int, str string) string {
	//var TextColor = []int{TextBlack, TextRed, TextGreen, TextYellow, TextBlue, TextCyan, TextWhite}
	//for _, v := range TextColor {
	//	if color == v {
	//		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
	//	}
	//}
	//return str

	switch color {
	case TextBlack:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlack, str)
	case TextRed:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextRed, str)
	case TextGreen:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextGreen, str)
	case TextYellow:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextYellow, str)
	case TextBlue:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlue, str)
	case TextMagenta:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextMagenta, str)
	case TextCyan:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextCyan, str)
	case TextWhite:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextWhite, str)
	default:
		return str
	}
}
