// Package render turns domain values into Discord message content: emoji, flags
// and Markdown. The domain packages hold no presentation.
package render

import (
	"fmt"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/timestamp"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

// RandomEmoji marks a choice the generator will fill in.
const RandomEmoji = "🎲"

const randomLabel = RandomEmoji + " Random"

// Challenge renders a generated challenge as Discord message content.
func Challenge(c challenge.Model) string {
	s := c.Stage()

	return strings.Join([]string{
		Stage(s),
		WeatherDescription(s.Location(), c.Weather()),
		Car(c.Car()),
		"",
	}, "\n")
}

// CompletionsByUser renders each user's recorded times as a block of text,
// keyed by user ID.
func CompletionsByUser(c challenge.Model) map[string]string {
	byUser := c.UserCompletions()

	out := make(map[string]string, len(byUser))
	for userID, durations := range byUser {
		buf := strings.Builder{}
		for _, d := range durations {
			buf.WriteString(timestamp.Format(d) + "\n")
		}
		out[userID] = buf.String()
	}

	return out
}

// StageConfig renders the stage half of an in-progress challenge builder.
func StageConfig(c challenge.Config) string {
	locationString := randomLabel
	stageString := randomLabel
	weatherString := randomLabel

	if c.Location != nil {
		locationString = Flag(*c.Location) + " " + c.Location.String()
	}

	switch {
	case c.Stage != nil:
		stageString = DistanceEmoji(c.Stage.Distance()) + " " + c.Stage.Name()
	case c.Distance != nil:
		stageString = DistanceEmoji(*c.Distance) + " Random " + c.Distance.String()
	}

	locationHasOneWeatherType := c.Location != nil && len(c.Location.Weather()) == 1

	switch {
	case c.Weather != nil:
		weatherString = WeatherEmoji(*c.Weather) + " " + c.Weather.String()
	case locationHasOneWeatherType:
		weatherString = fmt.Sprintf("%s *(probably %s though)*", randomLabel, c.Location.Weather()[0].String())
	}

	return fmt.Sprintf("Location: %s\nStage: %s\nWeather: %s", locationString, stageString, weatherString)
}

// CarConfig renders the car half of an in-progress challenge builder.
func CarConfig(c challenge.Config) string {
	drivetrainString := randomLabel
	classString := randomLabel
	carString := randomLabel

	switch {
	case c.Car != nil:
		drivetrainString = Drivetrain(c.Car.Class().Drivetrain())
		classString = c.Car.Class().String()
		carString = c.Car.Name()
	case c.Class != nil:
		drivetrainString = Drivetrain(c.Class.Drivetrain())
		classString = c.Class.String()
	case c.Drivetrain != nil:
		drivetrainString = Drivetrain(*c.Drivetrain)
	}

	return fmt.Sprintf("Drivetrain: %s\nClass: %s\nCar: %s", drivetrainString, classString, carString)
}

// Stage renders a stage as a bold "location » name" line.
func Stage(s stage.Model) string {
	return fmt.Sprintf("%s **%s » %s**", Flag(s.Location()), s.Location().String(), s.Name())
}

// Car renders a car as a bold "class » name" line.
func Car(c car.Model) string {
	return fmt.Sprintf("🏎️ **%s » %s**", c.Class().String(), c.Name())
}

// Drivetrain renders a drivetrain as its emoji followed by its name.
func Drivetrain(d drivetrain.Model) string {
	return DrivetrainEmoji(d) + " " + d.String()
}

// DistanceEmoji is the keycap emoji for a stage distance.
func DistanceEmoji(d stage.Distance) string {
	switch d {
	case stage.Short:
		return "4️⃣"
	case stage.Long:
		return "8️⃣"
	case stage.ReallyLong:
		return "♾️"
	case stage.Unknown:
		return "❓"
	}
	return "invalid distance"
}

// DrivetrainEmoji is the vehicle emoji for a drivetrain.
func DrivetrainEmoji(d drivetrain.Model) string {
	switch d {
	case drivetrain.FWD:
		return "🚗"
	case drivetrain.AWD:
		return "🚙"
	case drivetrain.AWDHYBRID:
		return "⚡"
	case drivetrain.RWD:
		return "🏎️"
	}
	return "invalid drivetrain"
}

// WeatherEmoji is the emoji for a weather type.
func WeatherEmoji(w weather.Model) string {
	switch w {
	case weather.DRY:
		return "☀️"
	case weather.WET:
		return "💧"
	case weather.SNOW:
		return "❄️"
	}
	return "invalid weather"
}
