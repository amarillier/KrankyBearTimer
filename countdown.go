package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-vgo/robotgo/base"
	_ "github.com/go-vgo/robotgo/key"
	_ "github.com/go-vgo/robotgo/mouse"
	_ "github.com/go-vgo/robotgo/screen"
	_ "github.com/go-vgo/robotgo/window"
)

// isValidDate checks if the entered date is valid in YYYY-MM-DD format
func isValidDate(d string) bool {
	_, err := time.Parse("2006-01-02", d)
	return err == nil
}

// selectDateDialog creates a dialog for selecting a date using Fyne calendar widget
func selectDateDialog(a fyne.App, parent fyne.Window, label string, currentDate string, callback func(string)) {
	var datePickerWindow fyne.Window

	// Parse current date or use today as default
	var selectedTime time.Time
	now := time.Now()
	if currentDate != "" && isValidDate(currentDate) {
		parsed, err := time.Parse("2006-01-02", currentDate)
		if err == nil {
			selectedTime = parsed
		} else {
			selectedTime = now
		}
	} else {
		selectedTime = now
	}

	datePickerWindow = a.NewWindow(fmt.Sprintf("Select Date: %s", label))
	_, monthNow, _ := time.Now().Date()
	if monthNow == time.December {
		datePickerWindow.SetIcon(resourceKrankyBearChristmasGrinchPng)
	} else {
		datePickerWindow.SetIcon(resourceKrankyBearBeretPng)
	}
	datePickerWindow.Resize(fyne.NewSize(350, 350))

	// Create calendar widget
	calendar := widget.NewCalendar(selectedTime, func(selectedDate time.Time) {
		dateStr := selectedDate.Format("2006-01-02")
		datePickerWindow.Close()
		datePickerWindow = nil
		callback(dateStr)
	})

	clearButton := widget.NewButton("Clear", func() {
		datePickerWindow.Close()
		datePickerWindow = nil
		callback("")
	})

	cancelButton := widget.NewButton("Cancel", func() {
		datePickerWindow.Close()
		datePickerWindow = nil
	})

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Select Date: %s", label)),
		widget.NewSeparator(),
		calendar,
		widget.NewSeparator(),
		container.NewHBox(clearButton, cancelButton),
	)

	datePickerWindow.SetContent(content)
	datePickerWindow.SetCloseIntercept(func() {
		datePickerWindow.Close()
		datePickerWindow = nil
	})
	datePickerWindow.Show()
}

func makeCountdownDates(a fyne.App, w fyne.Window, bg fyne.Canvas) {
	if countdown != nil && countdown.Content().Visible() {
		countdown.Show()
		return
	}
	if countdown == nil || !countdown.Content().Visible() {
		if countdown == nil {
			countdown = a.NewWindow(appName + ": Countdown Dates")
		}
		_, month, _ := time.Now().Date()
		if month == time.December {
			countdown.SetIcon(resourceKrankyBearChristmasGrinchPng)
		} else {
			countdown.SetIcon(resourceKrankyBearBeretPng)
		}

		// Get colors for the window
		var bre, bgr, bbl, ba uint8
		colors := strings.Split(bgcolor, ",")
		if len(colors) == 4 {
			col, _ := strconv.ParseUint(colors[0], 10, 8)
			bre = uint8(col)
			col, _ = strconv.ParseUint(colors[1], 10, 8)
			bgr = uint8(col)
			col, _ = strconv.ParseUint(colors[2], 10, 8)
			bbl = uint8(col)
			col, _ = strconv.ParseUint(colors[3], 10, 8)
			ba = uint8(col)
		}

		var tre, tgr, tbl, ta uint8
		colors = strings.Split(timecolor, ",")
		if len(colors) == 4 {
			col, _ := strconv.ParseUint(colors[0], 10, 8)
			tre = uint8(col)
			col, _ = strconv.ParseUint(colors[1], 10, 8)
			tgr = uint8(col)
			col, _ = strconv.ParseUint(colors[2], 10, 8)
			tbl = uint8(col)
			col, _ = strconv.ParseUint(colors[3], 10, 8)
			ta = uint8(col)
		}

		var dre, dgr, dbl, da uint8
		colors = strings.Split(datecolor, ",")
		if len(colors) == 4 {
			col, _ := strconv.ParseUint(colors[0], 10, 8)
			dre = uint8(col)
			col, _ = strconv.ParseUint(colors[1], 10, 8)
			dgr = uint8(col)
			col, _ = strconv.ParseUint(colors[2], 10, 8)
			dbl = uint8(col)
			col, _ = strconv.ParseUint(colors[3], 10, 8)
			da = uint8(col)
		}

		bgColorRect := color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
		background := canvas.NewRectangle(bgColorRect)

		// Helper function to edit label/description
		editLabelDialog := func(a fyne.App, parent fyne.Window, currentLabel string, callback func(string)) {
			var labelWindow fyne.Window

			labelWindow = a.NewWindow("Edit Label/Description")
			_, monthNow, _ := time.Now().Date()
			if monthNow == time.December {
				labelWindow.SetIcon(resourceKrankyBearChristmasGrinchPng)
			} else {
				labelWindow.SetIcon(resourceKrankyBearBeretPng)
			}
			labelWindow.Resize(fyne.NewSize(350, 150))

			labelEntry := widget.NewEntry()
			labelEntry.SetPlaceHolder("Enter label/description (max 20 chars)")
			if currentLabel != "" {
				if len(currentLabel) > 20 {
					currentLabel = currentLabel[:20]
				}
				labelEntry.SetText(currentLabel)
			}

			// Limit to 20 characters
			labelEntry.OnChanged = func(s string) {
				if len(s) > 20 {
					labelEntry.SetText(s[:20])
				}
			}

			messageLabel := widget.NewLabel("Max 20 characters")

			setButton := widget.NewButton("Set", func() {
				labelText := labelEntry.Text
				if len(labelText) > 20 {
					labelText = labelText[:20]
				}
				labelWindow.Close()
				labelWindow = nil
				callback(labelText)
			})

			cancelButton := widget.NewButton("Cancel", func() {
				labelWindow.Close()
				labelWindow = nil
			})

			content := container.NewVBox(
				widget.NewLabel("Label/Description:"),
				labelEntry,
				messageLabel,
				container.NewHBox(setButton, cancelButton),
			)

			labelWindow.SetContent(content)
			labelWindow.SetCloseIntercept(func() {
				labelWindow.Close()
				labelWindow = nil
			})
			labelWindow.Show()
		}

		// Helper function to create a date entry row
		createDateRow := func(num int, dateVar *string, descVar *string, daysText *canvas.Text) *fyne.Container {
			// Get label text - use description if available, otherwise default
			getLabelText := func() string {
				if *descVar != "" {
					return *descVar
				}
				return fmt.Sprintf("Date %d", num)
			}

			// Label button that opens dialog to edit description
			labelButtonText := getLabelText()
			if len(labelButtonText) > 20 {
				labelButtonText = labelButtonText[:20]
			}
			// Set a fixed pixel width that accommodates 20 wide characters (like W or M)
			// Calculate based on font size, but keep it compact to fit on one row
			fixedButtonWidth := float32(150) // Base width in pixels for ~20 characters (compact for single row)
			if datesize > 20 {
				// Scale with font size: datesize * 7 gives good width for wide chars (compact)
				fixedButtonWidth = float32(datesize) * 7
			}
			// Ensure it doesn't get too wide - keep it reasonable for single-row layout
			if fixedButtonWidth > 200 {
				fixedButtonWidth = 200
			}

			// Create label button - use regular button but enforce fixed width
			labelButton := widget.NewButton(labelButtonText, nil)
			testButton := widget.NewButton("Test", nil)
			buttonHeight := testButton.MinSize().Height

			// Force the button to have a fixed width using WithoutLayout container
			// This prevents Fyne from auto-sizing based on text content
			labelButtonContainer := container.NewWithoutLayout()
			labelButtonContainer.Resize(fyne.NewSize(fixedButtonWidth, buttonHeight))
			labelButtonContainer.Add(labelButton)
			labelButton.Resize(fyne.NewSize(fixedButtonWidth, buttonHeight))
			labelButton.Move(fyne.NewPos(0, 0))

			// Update label button text - maintain fixed width
			updateLabelButton := func() {
				newLabel := getLabelText()
				if len(newLabel) > 20 {
					newLabel = newLabel[:20]
				}
				labelButton.SetText(newLabel)
				// CRITICAL: Re-apply fixed size after text change to prevent auto-resize
				labelButton.Resize(fyne.NewSize(fixedButtonWidth, buttonHeight))
				labelButton.Move(fyne.NewPos(0, 0))
				labelButtonContainer.Resize(fyne.NewSize(fixedButtonWidth, buttonHeight))
				labelButton.Refresh()
				labelButtonContainer.Refresh()
			}

			// Edit label function
			editLabel := func() {
				editLabelDialog(a, countdown, *descVar, func(newLabel string) {
					*descVar = newLabel
					a.Preferences().SetString(fmt.Sprintf("countdown.desc%d", num), *descVar)
					updateLabelButton()
				})
			}

			labelButton.OnTapped = editLabel

			// Date button to trigger date picker
			dateButtonText := "Select Date"
			if *dateVar != "" {
				dateButtonText = *dateVar
			}
			dateButton := widget.NewButton(dateButtonText, nil)

			// Calculate fixed width for date button - needs to fit both "Select Date" and date strings (YYYY-MM-DD)
			// "Select Date" is the longest, so base width on that
			testDateButton := widget.NewButton("Select Date", nil)
			fixedDateButtonWidth := testDateButton.MinSize().Width
			// Add some padding to ensure it looks good
			fixedDateButtonWidth += 20

			// Create date button container with fixed width for alignment
			dateButtonContainer := container.NewWithoutLayout()
			dateButtonContainer.Resize(fyne.NewSize(fixedDateButtonWidth, buttonHeight))
			dateButtonContainer.Add(dateButton)
			dateButton.Resize(fyne.NewSize(fixedDateButtonWidth, buttonHeight))
			dateButton.Move(fyne.NewPos(0, 0))

			// Update function
			updateDays := func(dateStr string) {
				if dateStr != "" && isValidDate(dateStr) {
					days, err := daysUntil(dateStr)
					if err == nil {
						if days < 0 {
							daysText.Text = fmt.Sprintf("%d days ago", -days)
							daysText.Color = color.RGBA{R: 255, G: 100, B: 100, A: 255} // Red for past dates
						} else if days == 0 {
							daysText.Text = "Today!"
							daysText.Color = color.RGBA{R: 255, G: 215, B: 0, A: 255} // Gold for today
						} else if days == 1 {
							daysText.Text = "Tomorrow!"
							daysText.Color = color.RGBA{R: 100, G: 255, B: 100, A: 255} // Green for tomorrow
						} else {
							daysText.Text = fmt.Sprintf("%d days", days)
							daysText.Color = color.RGBA{R: tre, G: tgr, B: tbl, A: ta} // Use title color (will be updated to date color by updateCountdownColors)
						}
						*dateVar = dateStr
						dateButton.SetText(dateStr)
						// Re-apply fixed width after text change
						dateButton.Resize(fyne.NewSize(fixedDateButtonWidth, buttonHeight))
						dateButton.Move(fyne.NewPos(0, 0))
						dateButtonContainer.Resize(fyne.NewSize(fixedDateButtonWidth, buttonHeight))
						a.Preferences().SetString(fmt.Sprintf("countdown.date%d", num), *dateVar)
					} else {
						daysText.Text = "Invalid date"
						daysText.Color = color.RGBA{R: 255, G: 100, B: 100, A: 255}
					}
				} else if dateStr == "" {
					daysText.Text = ""
					*dateVar = ""
					dateButton.SetText("Select Date")
					// Re-apply fixed width after text change
					dateButton.Resize(fyne.NewSize(fixedDateButtonWidth, buttonHeight))
					dateButton.Move(fyne.NewPos(0, 0))
					dateButtonContainer.Resize(fyne.NewSize(fixedDateButtonWidth, buttonHeight))
					a.Preferences().SetString(fmt.Sprintf("countdown.date%d", num), "")
				}
				daysText.Refresh()
			}

			// Date picker function
			selectDate := func() {
				selectDateDialog(a, countdown, getLabelText(), *dateVar, func(selectedDate string) {
					updateDays(selectedDate)
				})
			}

			dateButton.OnTapped = selectDate

			daysText.TextSize = float32(datesize) * 0.8 // Same size as help text
			daysText.TextStyle = fyne.TextStyle{Bold: true}

			// Set initial days if date exists
			if *dateVar != "" {
				updateDays(*dateVar)
			}

			// Clear button for both date and label
			clearButton := widget.NewButton("Clear", func() {
				updateDays("")
				// updateDays already handles setting "Select Date" and resizing, but ensure it's done
				dateButton.SetText("Select Date")
				dateButton.Resize(fyne.NewSize(fixedDateButtonWidth, buttonHeight))
				dateButton.Move(fyne.NewPos(0, 0))
				dateButtonContainer.Resize(fyne.NewSize(fixedDateButtonWidth, buttonHeight))
				*descVar = ""
				a.Preferences().SetString(fmt.Sprintf("countdown.desc%d", num), "")
				updateLabelButton()
			})

			// Create HBox with proper spacing - ensure label button is on the left
			// Use WithoutLayout to prevent any wrapping and force single-row layout
			spacer := widget.NewLabel(" ") // Minimal spacer for visual separation

			clearButtonMinSize := clearButton.MinSize()

			// Create individual containers for each element to control sizing

			clearButtonContainer := container.NewWithoutLayout()
			clearButtonContainer.Resize(fyne.NewSize(clearButtonMinSize.Width, buttonHeight))
			clearButtonContainer.Add(clearButton)
			clearButton.Resize(fyne.NewSize(clearButtonMinSize.Width, buttonHeight))
			clearButton.Move(fyne.NewPos(0, 0))

			// Position and size elements explicitly in WithoutLayout container
			rowContainer := container.NewWithoutLayout()
			rowContainer.Resize(fyne.NewSize(700, buttonHeight)) // Fixed row width

			// Position elements horizontally, all at y=0 for same row alignment
			xPos := float32(0)

			// Label button container
			labelButtonContainer.Resize(fyne.NewSize(fixedButtonWidth, buttonHeight))
			labelButtonContainer.Move(fyne.NewPos(xPos, 0))
			rowContainer.Add(labelButtonContainer)
			xPos += fixedButtonWidth + 5 // 5px spacing

			// Spacer
			spacer.Resize(fyne.NewSize(5, buttonHeight))
			spacer.Move(fyne.NewPos(xPos, 0))
			rowContainer.Add(spacer)
			xPos += 5

			// Date button container
			dateButtonContainer.Move(fyne.NewPos(xPos, 0))
			rowContainer.Add(dateButtonContainer)
			xPos += fixedDateButtonWidth + 5

			// Clear button container
			clearButtonContainer.Move(fyne.NewPos(xPos, 0))
			rowContainer.Add(clearButtonContainer)
			xPos += clearButtonMinSize.Width + 5

			// Days text goes at the end
			daysTextContainer := container.NewWithoutLayout()
			daysTextMinSize := daysText.MinSize()
			daysTextContainer.Resize(fyne.NewSize(daysTextMinSize.Width, buttonHeight))
			daysTextContainer.Add(daysText)
			daysText.Move(fyne.NewPos(0, 0))
			daysTextContainer.Move(fyne.NewPos(xPos, 0))
			rowContainer.Add(daysTextContainer)

			dateRow := rowContainer

			return container.NewVBox(dateRow)
		}

		// Create three date rows
		// Initialize with time color, will be updated by updateDays function and updateCountdownColors
		daysText1 := canvas.NewText("", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
		daysText2 := canvas.NewText("", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
		daysText3 := canvas.NewText("", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})

		// Store in global variables for dynamic updates
		countdownDaysText1 = daysText1
		countdownDaysText2 = daysText2
		countdownDaysText3 = daysText3

		row1 := createDateRow(1, &countdownDate1, &countdownDesc1, daysText1)
		row2 := createDateRow(2, &countdownDate2, &countdownDesc2, daysText2)
		row3 := createDateRow(3, &countdownDate3, &countdownDesc3, daysText3)

		titleText := canvas.NewText("Days Until Countdown", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
		titleText.TextSize = float32(timesize)
		titleText.TextStyle = fyne.TextStyle{Bold: true}
		titleText.Alignment = fyne.TextAlignCenter
		countdownTitleText = titleText // Store for dynamic updates

		helpText := canvas.NewText("Click label button to edit name, or 'Select Date' for calendar. Clear removes both.", color.RGBA{R: dre, G: dgr, B: dbl, A: da})
		helpText.TextSize = float32(datesize) * 0.8
		helpText.Alignment = fyne.TextAlignCenter
		countdownHelpText = helpText // Store for dynamic updates

		countdownBackground = background // Store for dynamic updates

		closeButton := widget.NewButton("Close", func() {
			countdown.Close()
			countdown = nil
			// Clear global references when window is closed
			countdownTitleText = nil
			countdownHelpText = nil
			countdownBackground = nil
			countdownDaysText1 = nil
			countdownDaysText2 = nil
			countdownDaysText3 = nil
		})
		closeButton.Importance = widget.WarningImportance

		content := container.NewVBox(
			titleText,
			helpText,
			widget.NewSeparator(),
			row1,
			widget.NewSeparator(),
			row2,
			widget.NewSeparator(),
			row3,
			widget.NewSeparator(),
			container.NewCenter(closeButton),
		)

		mainContent := container.NewStack(background, content)
		countdown.SetContent(mainContent)
		countdown.Resize(fyne.NewSize(800, 350)) // Wider window to accommodate all elements on one row
		countdown.SetCloseIntercept(func() {
			countdown.Close()
			countdown = nil
			// Clear global references when window is closed
			countdownTitleText = nil
			countdownHelpText = nil
			countdownBackground = nil
			countdownDaysText1 = nil
			countdownDaysText2 = nil
			countdownDaysText3 = nil
		})
		countdown.Show()
		
		// Update countdown colors to ensure they're in sync with clock settings
		updateCountdownColors()
	}
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
