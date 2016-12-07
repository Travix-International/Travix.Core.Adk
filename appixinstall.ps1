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

Function Test-RegistryKeyValue
{
    # see: http://stackoverflow.com/questions/5648931/test-if-registry-value-exists
    param(
        # The path to the registry key where the value should be set.  Will be created if it doesn't exist.
        [Parameter(Mandatory=$true)]
        [string]
        $Path,

        # The name of the value being set.
        [Parameter(Mandatory=$true)]
        [string]
        $Name
    )

    if( -not (Test-Path -Path $Path -PathType Container) )
    {
        return $false
    }

    $properties = Get-ItemProperty -Path $Path 
    if( -not $properties )
    {
        return $false
    }

    $member = Get-Member -InputObject $properties -Name $Name
    if( $member )
    {
        return $true
    }
    else
    {
        return $false
    }

}

Function AddTo-SystemPath {
Param(
  [string]$Path
  )
  $registryPath = "Registry::HKEY_CURRENT_USER\Environment"

  if (Test-RegistryKeyValue -Path $registryPath -Name PATH)
  {
    $oldpath = (Get-ItemProperty -Path $registryPath -Name PATH).path
  }
  else
  {
    $oldpath = ""
  }

  #if($oldpath -Match $Path) {
  if($oldpath.Contains($Path)) {
    Write-Output "The folder is already in the PATH."
    return
  }

  # If we have an empty string, or if we already end with a semicolon, then append the path
  if($oldpath.EndsWith(";") -or $oldpath.Length -eq 0) {
    $newpath = "$oldpath$Path"
  }
  else {
    $newpath = "$oldpath;$Path"
  }

  Set-ItemProperty -Path $registryPath -Name PATH -Value $newPath

  # Updating the path for the current session
  $env:Path = $newpath

  # Broadcasting a settings change event so it's picked up by the other processes
  Broadcast-WMSettingsChange
}

Write-Output "Starting the Appix ADK installation"

$appixVersion = $env:APPIX_VERSION
if ([string]::IsNullOrEmpty($appixVersion)){
  # Determine the latest version
  $req = [System.Net.WebRequest]::Create("https://github.com/Travix-International/Travix.Core.Adk/releases/latest") -as [System.Net.HttpWebRequest]
  $req.Accept = "application/json"
  $res = $req.GetResponse()
  $outputStream = $res.GetResponseStream()
  $reader = New-Object System.IO.StreamReader $outputStream
  $content = $reader.ReadToEnd()

  # The releases are returned in a json like {... "tag_name":"hello-1.0.0.11", ...}, we have to extract the tag_name.
  $json = $content | ConvertFrom-Json
  $latestVersion = $json.tag_name
  $url = "https://github.com/Travix-International/Travix.Core.Adk/releases/download/$latestVersion/appix.exe"
}
else {
  # The version was explicitly specified
  $url = "https://github.com/Travix-International/Travix.Core.Adk/releases/download/$appixVersion/appix.exe"
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
