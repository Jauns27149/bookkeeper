### 打包

```bash
fyne-cross android -app-id com.janus.bookkeeper
adb install -r fyne-cross/dist/android/bookkeeper.apk

run-as com.janus.bookkeeper

# 备份
adb exec-out run-as com.janus.bookkeeper cat files/fyne/preferences.json > preferences.json
```


# x11环境
```bash
sudo pacman -S xorg-server xorg-xinit xorg-apps
```
