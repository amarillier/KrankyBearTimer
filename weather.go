package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// WeatherData holds the weather information
type WeatherData struct {
	CurrentTemp   float64
	MinTemp       float64
	MaxTemp       float64
	Precipitation float64
	Snowfall      float64
	Conditions    string
	LastUpdated   time.Time
	LocationName  string
}

var weatherWindow fyne.Window
var weatherBackground *canvas.Rectangle
var weatherTexts []*canvas.Text
var weatherEnabled bool
var weatherLatitude float64 = 0.0
var weatherLongitude float64 = 0.0
var currentWeatherData *WeatherData
var weatherRefreshTicker *time.Ticker
var weatherLocationName string // Store the city/country name

// loadWeatherSettings loads weather preferences
func loadWeatherSettings(a fyne.App) {
	weatherEnabled = a.Preferences().BoolWithFallback("weather.enabled", false)
	weatherLatitude = float64(a.Preferences().FloatWithFallback("weather.latitude", 0.0))
	weatherLongitude = float64(a.Preferences().FloatWithFallback("weather.longitude", 0.0))
	weatherLocationName = a.Preferences().StringWithFallback("weather.locationname", "")

	// Default to a reasonable location if not set (e.g., New York City)
	if weatherLatitude == 0.0 && weatherLongitude == 0.0 {
		weatherLatitude = 40.7128
		weatherLongitude = -74.0060
		weatherLocationName = "New York, United States"
	}
}

// saveWeatherSettings saves weather preferences
func saveWeatherSettings(a fyne.App) {
	a.Preferences().SetBool("weather.enabled", weatherEnabled)
	a.Preferences().SetFloat("weather.latitude", weatherLatitude)
	a.Preferences().SetFloat("weather.longitude", weatherLongitude)
	a.Preferences().SetString("weather.locationname", weatherLocationName)
}

// searchLocation searches for a city/country using Open-Meteo Geocoding API
func searchLocation(query string) (name string, lat, lon float64, err error) {
	if query == "" {
		return "", 0, 0, fmt.Errorf("search query is empty")
	}

	// Open-Meteo Geocoding API
	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json",
		strings.ReplaceAll(query, " ", "%20"))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to search location: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, 0, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to read response: %v", err)
	}

	var apiResponse struct {
		Results []struct {
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Country   string  `json:"country"` // Full country name
			Admin1    string  `json:"admin1"`  // State/Province
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", 0, 0, fmt.Errorf("failed to parse geocoding data: %v", err)
	}

	if len(apiResponse.Results) == 0 {
		return "", 0, 0, fmt.Errorf("no locations found for: %s", query)
	}

	result := apiResponse.Results[0]
	// Format location name: "City, State, Country" or "City, Country"
	name = result.Name
	if result.Admin1 != "" {
		name += ", " + result.Admin1
	}
	if result.Country != "" {
		name += ", " + result.Country
	}

	return name, result.Latitude, result.Longitude, nil
}

// fetchWeather fetches weather data from Open-Meteo API
func fetchWeather(a fyne.App) (*WeatherData, error) {
	if weatherLatitude == 0.0 && weatherLongitude == 0.0 {
		return nil, fmt.Errorf("location not set")
	}

	// Open-Meteo API endpoint for current weather and forecast
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,precipitation,snowfall&daily=temperature_2m_max,temperature_2m_min,precipitation_sum,snowfall_sum&timezone=auto",
		weatherLatitude, weatherLongitude)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var apiResponse struct {
		Current struct {
			Temperature2m float64 `json:"temperature_2m"`
			Precipitation float64 `json:"precipitation"`
			Snowfall      float64 `json:"snowfall"`
			Time          string  `json:"time"`
		} `json:"current"`
		Daily struct {
			Time             []string  `json:"time"`
			Temperature2mMax []float64 `json:"temperature_2m_max"`
			Temperature2mMin []float64 `json:"temperature_2m_min"`
			PrecipitationSum []float64 `json:"precipitation_sum"`
			SnowfallSum      []float64 `json:"snowfall_sum"`
		} `json:"daily"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		bodyPreview := string(body)
		if len(bodyPreview) > 200 {
			bodyPreview = bodyPreview[:200] + "..."
		}
		return nil, fmt.Errorf("failed to parse weather data: %v (response: %s)", err, bodyPreview)
	}

	weather := &WeatherData{
		CurrentTemp:   apiResponse.Current.Temperature2m,
		LastUpdated:   time.Now(),
		Precipitation: apiResponse.Current.Precipitation,
		Snowfall:      apiResponse.Current.Snowfall,
	}

	// Get today's forecast (index 0)
	if len(apiResponse.Daily.Temperature2mMax) > 0 {
		weather.MaxTemp = apiResponse.Daily.Temperature2mMax[0]
	}
	if len(apiResponse.Daily.Temperature2mMin) > 0 {
		weather.MinTemp = apiResponse.Daily.Temperature2mMin[0]
	}

	// Determine conditions
	var conditions []string
	if weather.Precipitation > 0.1 {
		conditions = append(conditions, "Rain expected")
	}
	if weather.Snowfall > 0.1 {
		conditions = append(conditions, "Snow expected")
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "Clear")
	}
	weather.Conditions = strings.Join(conditions, ", ")

	// Use stored location name if available, otherwise show coordinates
	if weatherLocationName != "" {
		weather.LocationName = weatherLocationName
	} else {
		weather.LocationName = fmt.Sprintf("Lat: %.2f, Lon: %.2f", weatherLatitude, weatherLongitude)
	}

	currentWeatherData = weather
	return weather, nil
}

// updateWeatherDisplay updates the weather window with current data
func updateWeatherDisplay(a fyne.App) {
	// Always try to fetch weather when enabled (even if window is closed) to update clock display
	if !weatherEnabled {
		return
	}

	weather, err := fetchWeather(a)
	if err != nil {
		// If window is open, show error there
		if weatherWindow != nil {
			if len(weatherTexts) > 0 {
				fyne.Do(func() {
					if weatherTexts[0] != nil {
						weatherTexts[0].Text = fmt.Sprintf("Error: %v", err)
						weatherTexts[0].Refresh()
					}
					// Clear other fields on error
					for i := 1; i < len(weatherTexts) && i < 5; i++ {
						if weatherTexts[i] != nil {
							weatherTexts[i].Text = ""
							weatherTexts[i].Refresh()
						}
					}
				})
			}
		}
		// currentWeatherData remains unchanged on error, so clock display keeps showing last known data
		return
	}
	if weather == nil {
		return
	}

	// currentWeatherData is updated by fetchWeather, so clock display will update automatically
	// Only update window display if window is open
	if weatherWindow == nil {
		return
	}

	// Helper function to convert Celsius to Fahrenheit
	celsiusToFahrenheit := func(c float64) float64 {
		return c*9/5 + 32
	}

	// Update text elements
	fyne.Do(func() {
		if len(weatherTexts) >= 5 {
			if weatherTexts[0] != nil {
				weatherTexts[0].Text = weather.LocationName
				weatherTexts[0].Refresh()
			}
			if weatherTexts[1] != nil {
				// Show current temp with Fahrenheit first, then Celsius
				weatherTexts[1].Text = fmt.Sprintf("Current: %.1f°F / %.1f°C", celsiusToFahrenheit(weather.CurrentTemp), weather.CurrentTemp)
				weatherTexts[1].Refresh()
			}
			if weatherTexts[2] != nil {
				// Show high temp with Fahrenheit first, then Celsius
				weatherTexts[2].Text = fmt.Sprintf("High: %.1f°F / %.1f°C", celsiusToFahrenheit(weather.MaxTemp), weather.MaxTemp)
				weatherTexts[2].Refresh()
			}
			if weatherTexts[3] != nil {
				// Show low temp with Fahrenheit first, then Celsius
				weatherTexts[3].Text = fmt.Sprintf("Low: %.1f°F / %.1f°C", celsiusToFahrenheit(weather.MinTemp), weather.MinTemp)
				weatherTexts[3].Refresh()
			}
			if weatherTexts[4] != nil {
				precipText := ""
				if weather.Precipitation > 0.1 {
					precipText = fmt.Sprintf("Rain: %.1f mm", weather.Precipitation)
				}
				if weather.Snowfall > 0.1 {
					if precipText != "" {
						precipText += " | "
					}
					// Snowfall is in cm from API
					precipText += fmt.Sprintf("Snow: %.1f cm", weather.Snowfall)
				}
				if precipText == "" {
					precipText = "No precipitation"
				}
				precipText += fmt.Sprintf(" | Updated: %s", weather.LastUpdated.Format("15:04:05"))
				weatherTexts[4].Text = precipText
				weatherTexts[4].Refresh()
			}
		}
	})
}

// startWeatherRefresh starts the background goroutine that refreshes weather every 5 minutes
func startWeatherRefresh(a fyne.App) {
	if !weatherEnabled {
		return
	}

	// Initial fetch - start immediately
	go func() {
		updateWeatherDisplay(a)
	}()

	// Start ticker for 5-minute refresh
	go func() {
		weatherRefreshTicker = time.NewTicker(5 * time.Minute)
		defer weatherRefreshTicker.Stop()

		for range weatherRefreshTicker.C {
			if weatherEnabled {
				updateWeatherDisplay(a)
			}
		}
	}()
}

// makeWeatherWindow creates the weather display window
func makeWeatherWindow(a fyne.App, w fyne.Window, bg fyne.Canvas) {
	if weatherWindow != nil && weatherWindow.Content().Visible() {
		weatherWindow.Show()
		weatherWindow.RequestFocus()
		return
	}

	if weatherWindow != nil {
		weatherWindow.Close()
		weatherWindow = nil
	}

	weatherWindow = a.NewWindow(appName + ": Weather")
	_, month, _ := time.Now().Date()
	if month == time.December {
		weatherWindow.SetIcon(resourceKrankyBearChristmasGrinchPng)
	} else {
		weatherWindow.SetIcon(resourceKrankyBearFedoraRedPng)
	}

	// Parse colors using helper function (from alarm.go)
	bre, bgr, bbl, ba := parseColor(bgcolor)
	tre, tgr, tbl, ta := parseColor(timecolor)

	// Background
	background := canvas.NewRectangle(color.RGBA{R: bre, G: bgr, B: bbl, A: ba})
	background.Resize(fyne.NewSize(400, 300))
	weatherBackground = background

	// Create text elements (reduced to 5 since high/low will be on same row)
	weatherTexts = make([]*canvas.Text, 5)
	weatherTexts[0] = canvas.NewText("Loading...", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
	weatherTexts[0].TextStyle = fyne.TextStyle{Bold: true}
	weatherTexts[0].Alignment = fyne.TextAlignCenter
	weatherTexts[0].TextSize = float32(utcsize)

	weatherTexts[1] = canvas.NewText("", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
	weatherTexts[1].TextStyle = fyne.TextStyle{Bold: true}
	weatherTexts[1].Alignment = fyne.TextAlignCenter
	weatherTexts[1].TextSize = float32(utcsize)

	weatherTexts[2] = canvas.NewText("", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
	weatherTexts[2].Alignment = fyne.TextAlignCenter
	weatherTexts[2].TextSize = float32(utcsize)

	weatherTexts[3] = canvas.NewText("", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
	weatherTexts[3].Alignment = fyne.TextAlignCenter
	weatherTexts[3].TextSize = float32(utcsize)

	weatherTexts[4] = canvas.NewText("", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
	weatherTexts[4].Alignment = fyne.TextAlignCenter
	weatherTexts[4].TextSize = float32(utcsize - 4)

	// Location settings - wider input fields
	latEntry := widget.NewEntry()
	latEntry.SetText(fmt.Sprintf("%.4f", weatherLatitude))
	latEntry.SetPlaceHolder("Latitude")
	latEntryContainer := container.NewWithoutLayout(latEntry)
	latEntryContainer.Resize(fyne.NewSize(150, latEntry.MinSize().Height))
	latEntry.Resize(fyne.NewSize(150, latEntry.MinSize().Height))

	lonEntry := widget.NewEntry()
	lonEntry.SetText(fmt.Sprintf("%.4f", weatherLongitude))
	lonEntry.SetPlaceHolder("Longitude")
	lonEntryContainer := container.NewWithoutLayout(lonEntry)
	lonEntryContainer.Resize(fyne.NewSize(150, lonEntry.MinSize().Height))
	lonEntry.Resize(fyne.NewSize(150, lonEntry.MinSize().Height))

	// City/Country search entry
	cityEntry := widget.NewEntry()
	if weatherLocationName != "" {
		cityEntry.SetText(weatherLocationName)
	} else {
		cityEntry.SetPlaceHolder("Enter city, country (e.g., Paris, France)")
	}
	// Wrap in container with fixed width to prevent button overlap
	cityEntryContainer := container.NewWithoutLayout(cityEntry)
	cityEntryContainer.Resize(fyne.NewSize(250, cityEntry.MinSize().Height))
	cityEntry.Resize(fyne.NewSize(250, cityEntry.MinSize().Height))
	cityEntry.Move(fyne.NewPos(0, 0))

	// Search button for city/country lookup
	searchButton := widget.NewButton("Search Location", func() {
		query := strings.TrimSpace(cityEntry.Text)
		if query == "" {
			return
		}

		// Show searching message
		if len(weatherTexts) > 0 && weatherTexts[0] != nil {
			fyne.Do(func() {
				weatherTexts[0].Text = "Searching..."
				weatherTexts[0].Refresh()
			})
		}

		// Search in background
		go func() {
			name, lat, lon, err := searchLocation(query)
			if err != nil {
				fyne.Do(func() {
					if len(weatherTexts) > 0 && weatherTexts[0] != nil {
						weatherTexts[0].Text = fmt.Sprintf("Error: %v", err)
						weatherTexts[0].Refresh()
					}
				})
				return
			}

			// Update location
			fyne.Do(func() {
				weatherLatitude = lat
				weatherLongitude = lon
				weatherLocationName = name
				latEntry.SetText(fmt.Sprintf("%.4f", lat))
				lonEntry.SetText(fmt.Sprintf("%.4f", lon))
				cityEntry.SetText(name)
				saveWeatherSettings(a)
				updateWeatherDisplay(a)
			})
		}()
	})

	enabledCheck := widget.NewCheck("Enable Weather", func(checked bool) {
		weatherEnabled = checked
		saveWeatherSettings(a)
		if checked {
			startWeatherRefresh(a)
			// Trigger immediate fetch for faster display
			go func() {
				updateWeatherDisplay(a)
			}()
		} else if weatherRefreshTicker != nil {
			weatherRefreshTicker.Stop()
			// Clear weather data when disabled
			currentWeatherData = nil
		}
		// Refresh clock display to show/hide weather text
		if clock != nil {
			updateTimezoneVisibility()
		}
	})
	enabledCheck.SetChecked(weatherEnabled)

	// Buttons
	updateButton := widget.NewButton("Update Location", func() {
		lat, err1 := strconv.ParseFloat(latEntry.Text, 64)
		lon, err2 := strconv.ParseFloat(lonEntry.Text, 64)
		if err1 == nil && err2 == nil {
			weatherLatitude = lat
			weatherLongitude = lon
			// Clear location name when manually updating coordinates
			weatherLocationName = ""
			cityEntry.SetText("")
			saveWeatherSettings(a)
			// Show loading message
			if len(weatherTexts) > 0 && weatherTexts[0] != nil {
				fyne.Do(func() {
					weatherTexts[0].Text = "Loading..."
					weatherTexts[0].Refresh()
				})
			}
			updateWeatherDisplay(a)
		}
	})

	refreshButton := widget.NewButton("Refresh Now", func() {
		// Show loading message
		if len(weatherTexts) > 0 && weatherTexts[0] != nil {
			fyne.Do(func() {
				weatherTexts[0].Text = "Loading..."
				weatherTexts[0].Refresh()
			})
		}
		updateWeatherDisplay(a)
	})

	// Layout with subtle spacing
	content := container.NewVBox()
	content.Add(weatherTexts[0]) // Location
	content.Add(weatherTexts[1]) // Current temp (with °C and °F)
	// High and Low on same row with spacing
	highLowRow := container.NewHBox(
		weatherTexts[2],
		widget.NewLabel("      "), // Spacer between high and low
		weatherTexts[3],
	)
	content.Add(highLowRow)
	content.Add(weatherTexts[4]) // Precipitation and last updated
	content.Add(widget.NewSeparator())
	// City/Country search - put button below aligned with input field
	searchLabel := widget.NewLabel("Search by City/Country:")
	content.Add(container.NewHBox(searchLabel, cityEntryContainer))
	// Align button directly under the input field (add padding to match label width)
	buttonPadding := container.NewHBox(
		widget.NewLabel("                       "), // Approximate padding to align with input field
		searchButton,
	)
	content.Add(buttonPadding)
	content.Add(widget.NewSeparator())
	// Separate Lat and Lon on different rows, aligned with consistent spacing
	// Use fixed-width label containers to ensure perfect alignment
	latLabel := widget.NewLabel("Latitude:")
	lonLabel := widget.NewLabel("Longitude:")

	// Calculate maximum label width to align both
	latWidth := latLabel.MinSize().Width
	lonWidth := lonLabel.MinSize().Width
	maxLabelWidth := latWidth
	if lonWidth > maxLabelWidth {
		maxLabelWidth = lonWidth
	}

	// Create fixed-width containers for labels
	latLabelContainer := container.NewWithoutLayout(latLabel)
	latLabelContainer.Resize(fyne.NewSize(maxLabelWidth, latLabel.MinSize().Height))
	latLabel.Move(fyne.NewPos(0, 0))

	lonLabelContainer := container.NewWithoutLayout(lonLabel)
	lonLabelContainer.Resize(fyne.NewSize(maxLabelWidth, lonLabel.MinSize().Height))
	lonLabel.Move(fyne.NewPos(0, 0))

	// Add consistent padding between label and input field (reduced for longitude since label is longer)
	latRow := container.NewHBox(latLabelContainer, widget.NewLabel("   "), latEntryContainer)
	lonRow := container.NewHBox(lonLabelContainer, widget.NewLabel(" "), lonEntryContainer) // Less padding for longitude
	content.Add(latRow)
	content.Add(lonRow)
	// Current location button next to coordinates
	currentLocationButton := widget.NewButton("Select Current Location", func() {
		// Show searching message
		if len(weatherTexts) > 0 && weatherTexts[0] != nil {
			fyne.Do(func() {
				weatherTexts[0].Text = "Detecting location..."
				weatherTexts[0].Refresh()
			})
		}

		// Get current location using IP geolocation (simplified approach)
		// Note: This is a basic implementation using a free IP geolocation service
		go func() {
			// Use ipapi.co free IP geolocation service
			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Get("https://ipapi.co/json/")
			if err != nil {
				fyne.Do(func() {
					if len(weatherTexts) > 0 && weatherTexts[0] != nil {
						weatherTexts[0].Text = fmt.Sprintf("Error getting location: %v", err)
						weatherTexts[0].Refresh()
					}
				})
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fyne.Do(func() {
					if len(weatherTexts) > 0 && weatherTexts[0] != nil {
						weatherTexts[0].Text = "Could not detect current location"
						weatherTexts[0].Refresh()
					}
				})
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fyne.Do(func() {
					if len(weatherTexts) > 0 && weatherTexts[0] != nil {
						weatherTexts[0].Text = fmt.Sprintf("Error reading location: %v", err)
						weatherTexts[0].Refresh()
					}
				})
				return
			}

			var location struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				City      string  `json:"city"`
				Region    string  `json:"region"`
				Country   string  `json:"country_name"`
			}

			if err := json.Unmarshal(body, &location); err != nil {
				fyne.Do(func() {
					if len(weatherTexts) > 0 && weatherTexts[0] != nil {
						weatherTexts[0].Text = fmt.Sprintf("Error parsing location: %v", err)
						weatherTexts[0].Refresh()
					}
				})
				return
			}

			// Update location
			fyne.Do(func() {
				weatherLatitude = location.Latitude
				weatherLongitude = location.Longitude
				// Format location name
				name := location.City
				if location.Region != "" {
					name += ", " + location.Region
				}
				if location.Country != "" {
					name += ", " + location.Country
				}
				weatherLocationName = name
				latEntry.SetText(fmt.Sprintf("%.4f", location.Latitude))
				lonEntry.SetText(fmt.Sprintf("%.4f", location.Longitude))
				cityEntry.SetText(name)
				saveWeatherSettings(a)
				updateWeatherDisplay(a)
			})
		}()
	})
	content.Add(container.NewHBox(currentLocationButton))
	content.Add(container.NewHBox(updateButton, refreshButton))
	content.Add(enabledCheck)

	scroll := container.NewScroll(content)
	padded := container.NewPadded(scroll)
	mainContent := container.NewStack(background, padded)
	weatherWindow.SetContent(mainContent)

	weatherWindow.SetCloseIntercept(func() {
		weatherWindow.Close()
		weatherWindow = nil
		weatherBackground = nil
		weatherTexts = nil
	})

	// Initial color update
	updateWeatherColors()

	// Set window size to minimum required and make it non-resizable
	weatherWindow.CenterOnScreen()
	// Let content determine minimum size, then set window to that size
	contentMinSize := content.MinSize()
	weatherWindow.Resize(fyne.NewSize(450, contentMinSize.Height+50)) // Add small padding for borders/padding
	weatherWindow.SetFixedSize(true)                                  // Make window non-resizable
	weatherWindow.Show()

	// Initial fetch - start immediately
	go func() {
		updateWeatherDisplay(a)
	}()
}

// updateWeatherColors updates the weather window colors
func updateWeatherColors() {
	if weatherWindow == nil {
		return
	}

	bre, bgr, bbl, ba := parseColor(bgcolor)
	tre, tgr, tbl, ta := parseColor(timecolor)

	fyne.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				weatherWindow = nil
				weatherBackground = nil
				weatherTexts = nil
			}
		}()

		if weatherBackground != nil {
			weatherBackground.FillColor = color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
			if weatherWindow != nil && weatherWindow.Content() != nil {
				contentSize := weatherWindow.Content().Size()
				weatherBackground.Resize(contentSize)
			}
			weatherBackground.Refresh()
		}

		for i := 0; i < len(weatherTexts); i++ {
			if weatherTexts[i] != nil {
				weatherTexts[i].Color = color.RGBA{R: tre, G: tgr, B: tbl, A: ta}
				weatherTexts[i].Refresh()
			}
		}

		if weatherWindow != nil && weatherWindow.Content() != nil {
			weatherWindow.Content().Refresh()
			if weatherWindow.Canvas() != nil {
				weatherWindow.Canvas().Refresh(weatherWindow.Content())
			}
		}
	})
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
