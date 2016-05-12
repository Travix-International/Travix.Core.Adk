function Download-File {
param (
  [string]$url,
  [string]$file
 )
  Write-Output "Downloading $url to $file"
  $downloader = new-object System.Net.WebClient

  $downloader.DownloadFile($url, $file)
}

Function Broadcast-WMSettingsChange {
  if (-not ("win32.nativemethods" -as [type])) {
    Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
   IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
   uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
  }

  $HWND_BROADCAST = [intptr]0xffff;
  $WM_SETTINGCHANGE = 0x1a;
  $result = [uintptr]::zero

  # notify all windows of environment block change
  [win32.nativemethods]::SendMessageTimeout($HWND_BROADCAST, $WM_SETTINGCHANGE, [uintptr]::Zero, "Environment", 2, 5000, [ref]$result) >$null 2>&1;
}

Function AddTo-SystemPath {
Param(
  [string]$Path
  )
  $oldpath = (Get-ItemProperty -Path "Registry::HKEY_CURRENT_USER\Environment" -Name PATH).path

  if($oldpath.EndsWith(";")) {
    $newpath = "$oldpath$Path"
  }
  else {
    $newpath = "$oldpath;$Path"
  }

  Set-ItemProperty -Path "Registry::HKEY_CURRENT_USER\Environment" -Name PATH -Value $newPath

  # Updating the path for the current session
  $env:Path = $newpath

  # Broadcasting a settings change event so it's picked up by the other processes
  Broadcast-WMSettingsChange
}

Write-Output "Starting the Appix ADK installation"

$url = "https://raw.githubusercontent.com/markvincze/jil-playground/master/appix-windows.exe"

$appixVersion = $env:appixVersion
if (![string]::IsNullOrEmpty($appixVersion)){
  # TODO: the file naming scheme needs to be finalized once we figure out how we build and publish the binaries
  $url = "https://raw.githubusercontent.com/markvincze/jil-playground/master/appix-windows-$appixVersion.exe"  
}

# We install into ~/.appix
$appixFolder = Join-Path "$env:USERPROFILE" ".appix"

if(!(Test-Path -Path $appixFolder)) {
  New-Item -Path $appixFolder -ItemType directory
}

$appixFile = Join-Path $appixFolder "appix.exe"

if(!(Test-Path -Path $appixFile)) {
  Write-Output "Downloading the appix binary to $appixFolder."
}
else {
  Write-Output "Appix is already installed in $appixFolder, trying to update."
}

Write-Output "Downloading the Appix ADK from $url."

Download-File $url $appixFile

Write-Output "Download completed. Adding appix folder to the PATH."

AddTo-SystemPath $appixFolder

Write-Output "The Appix ADK has been installed. You can start using it by typing appix."
Write-Output "(You might have to restart your terminal session to refresh your PATH.)"