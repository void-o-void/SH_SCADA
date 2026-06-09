import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import HMI.Device 1.0
import HMI.Core 1.0
import "../../antui" as Antui

ApplicationWindow {
    id: mainWindow
    width: 1000
    height: 600
    visible: true
    title: "SHlot平台"
    flags: Qt.Window | Qt.FramelessWindowHint
    color: "#00000000"

    FontLoader {
        id: titleFont
        source: "qrc:/font/YouSheBiaoTiHei-2.ttf"
    }

    FontLoader {
        id: font_deyihei
        source: "qrc:/font/SmileySans-Oblique.otf"
    }

    FontLoader {
        id: faBrands_Regular
        source: "qrc:/fontawesome6.72pro/Brands_Regular_400.otf"
    }

    property bool tabIsLogin: true

    ListModel {
        id: userHistory
    }

    Rectangle{
        radius: 20
        anchors.fill: parent
        clip: true
        Image {
            id: backgroundImage
            anchors.fill: parent
            source: "qrc:/image/login_bg4.png"
            fillMode: Image.PreserveAspectCrop
        }
    }

    Image {
        id: vxcode
        source: "qrc:/image/erweima.png"
        width: parent.width * 0.2
        height: parent.width * 0.2
        anchors {
            verticalCenter: parent.verticalCenter
            left: parent.left
            leftMargin: (parent.width/2-parent.width * 0.2)/2
        }
    }

    MouseArea {
        anchors {
            top: parent.top
            left: parent.left
            right: parent.right
        }
        height: 40
        property point clickPos

        onPressed: (mouse) => {
            clickPos = Qt.point(mouse.x, mouse.y)
        }

        onPositionChanged: (mouse) => {
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            mainWindow.x += delta.x
            mainWindow.y += delta.y
        }
    }

    Rectangle {
        color: "#F2FFFFFF"
        height: parent.height
        width: parent.width/2
        anchors.right: parent.right

        Button {
            id: linkButton
            text: "去注册账户？"
            visible: mainWindow.tabIsLogin
            anchors {
                right: parent.right
                top: parent.top
            }

            // 背景透明
            background: Rectangle {
                color: "transparent"
            }

            // 文字样式
            contentItem: Text {
                text: linkButton.text
                color: "#0066CC"  // 链接蓝色
                font {
                    pixelSize: 16
                    underline: true
                }
            }

            onClicked: {
                mainWindow.tabIsLogin = false
                console.log("切换到注册界面，tabIsLogin =", mainWindow.tabIsLogin)
            }
        }

        Text {
            id: login_title
            text: qsTr("SHiot平台")
            color: "#21C491"
            anchors {
                horizontalCenter: parent.horizontalCenter
                top: parent.top
                topMargin: parent.height * 0.14
            }
            font {
                family: titleFont.name
                pixelSize: 48
            }
        }

        Text {
            id: userInTip
            text: qsTr("请输入账号")
            visible: mainWindow.tabIsLogin
            color: "#1A1A1A"
            opacity: 0.3  // 添加透明度，0-1之间
            font {
                family: "Microsoft YaHei"
                pixelSize: 14
            }
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: login_title.bottom
                topMargin: 40
            }
        }

        Antui.FramelessCombox {
            id: userInBox
            visible: mainWindow.tabIsLogin
            editable: true
            width: parent.width * 0.56
            model: userHistory
            textRole: "name"
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: userInTip.bottom
            }

            font {
                family: "Microsoft YaHei"
                pixelSize: 18
            }
        }

        Text {
            id: pwdInTip
            text: qsTr("请输入密码")
            visible: mainWindow.tabIsLogin
            color: "#1A1A1A"
            opacity: 0.3  // 添加透明度，0-1之间
            font {
                family: "Microsoft YaHei"
                pixelSize: 14
            }
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: userInBox.bottom
                topMargin: 20
            }
        }

        Antui.FramelessLineEdit {
            id: pwdInBox
            visible: mainWindow.tabIsLogin
            width: parent.width * 0.56
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: pwdInTip.bottom
            }

            font {
                family: "Microsoft YaHei"
                pixelSize: 18
            }

            echoMode: TextInput.Password
        }

        RowLayout {
            id: savepwdselect
            visible: mainWindow.tabIsLogin
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: pwdInBox.bottom
                topMargin: 35
            }
            RadioButton {
                id: usrpwdselect
                implicitWidth: 15
                implicitHeight: 15

                // 隐藏默认的指示器和背景
                indicator: null

                // 自定义背景（使用SVG图标）
                background: Image {
                    anchors.fill: parent
                    source: usrpwdselect.checked ? "qrc:/icon/selected.svg" : "qrc:/icon/unselect.svg"
                    fillMode: Image.PreserveAspectFit
                    sourceSize: Qt.size(width, height)
                }

                // 可选：添加点击效果
                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: usrpwdselect.toggle()
                }
            }

            Text {
                text: qsTr("记住密码")
                font {
                    family: "Microsoft YaHei"
                    pixelSize: 14
                }
            }
        }

        RowLayout {
            visible: mainWindow.tabIsLogin
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: savepwdselect.bottom
                topMargin: 20
            }
            RadioButton {
                id: knowTipSelect
                implicitWidth: 15
                implicitHeight: 15

                // 隐藏默认的指示器和背景
                indicator: null

                // 自定义背景（使用SVG图标）
                background: Image {
                    anchors.fill: parent
                    source: knowTipSelect.checked ? "qrc:/icon/selected.svg" : "qrc:/icon/unselect.svg"
                    fillMode: Image.PreserveAspectFit
                    sourceSize: Qt.size(width, height)
                }

                // 可选：添加点击效果
                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: knowTipSelect.toggle()
                }
            }

            Text {
                text: qsTr("已阅读并同意隐私保护协议以及互联网……")
                font {
                    family: "Microsoft YaHei"
                    pixelSize: 14
                }
            }
        }

        Button {
            id: loginButton
            visible: mainWindow.tabIsLogin
            width: parent.width * 0.56
            height: 45
            text: "登录"

            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                bottom: parent.bottom
                bottomMargin: parent.height * 0.14
            }

            background: Rectangle {
                radius: 10  // 圆角
                color: loginButton.down ? "#3d8b40" :  "#5CBF60"

                border.width: 0
            }

            contentItem: Text {
                text: loginButton.text
                color: "white"  // 白色文字
                font {
                    pixelSize: 20
                    family: "Microsoft YaHei"
                    bold: true
                }
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
            }

            onClicked: {
                var username = userInBox.editText
                var password = pwdInBox.text
                if (UserMgr.login(username, password)) {
                    // 保存到登录历史（去重）
                    var found = false
                    for (var i = 0; i < userHistory.count; i++) {
                        if (userHistory.get(i).name === username) { found = true; break }
                    }
                    if (!found && username !== "") userHistory.append({"name": username})

                    var component = Qt.createComponent("MainWindow.qml")
                    if (component.status === Component.Ready) {
                        var mainWindow = component.createObject(null, {
                            "visible": true
                        })
                        mainWindow.visible = true
                    } else {
                        console.log("创建主窗口失败:", component.errorString())
                    }
                }
            }
        }

        //注册区域
        Button {
            id: gologinButton
            text: "已有账户，去登录？"
            visible: !mainWindow.tabIsLogin
            anchors {
                right: parent.right
                top: parent.top
            }

            // 背景透明
            background: Rectangle {
                color: "transparent"
            }

            // 文字样式
            contentItem: Text {
                text: gologinButton.text
                color: "#0066CC"  // 链接蓝色
                font {
                    pixelSize: 16
                    underline: true
                }
            }

            onClicked: {
                mainWindow.tabIsLogin = true
                console.log("切换到注册界面，tabIsLogin =", mainWindow.tabIsLogin)
            }
        }

        Text {
            id: regTitle
            visible: !mainWindow.tabIsLogin
            text: qsTr("让物联更有趣，欢迎加入！")
            font {
                family: font_deyihei.name
                pixelSize: 30
            }
            anchors {
                top: login_title.bottom
                topMargin: 15
                left: parent.left
                leftMargin: parent.width*0.22
            }
        }


        Text {
            id: regUsrTip
            text: qsTr("请输入账户")
            visible: !mainWindow.tabIsLogin
            color: "#1A1A1A"
            opacity: 0.3  // 添加透明度，0-1之间
            font {
                family: "Microsoft YaHei"
                pixelSize: 14
            }
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: login_title.bottom
                topMargin: 80
            }
        }

        Antui.FramelessLineEdit {
            id: regusrInBox
            visible: !mainWindow.tabIsLogin
            width: parent.width * 0.56
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: regUsrTip.bottom
            }

            font {
                family: "Microsoft YaHei"
                pixelSize: 18
            }
        }

        Text {
            id: regpwdTip
            text: qsTr("请输入密码")
            visible: !mainWindow.tabIsLogin
            color: "#1A1A1A"
            opacity: 0.3  // 添加透明度，0-1之间
            font {
                family: "Microsoft YaHei"
                pixelSize: 14
            }
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: regusrInBox.bottom
                topMargin: 20
            }
        }

        Antui.FramelessLineEdit {
            id: regPwdInBox
            visible: !mainWindow.tabIsLogin
            width: parent.width * 0.56
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: regpwdTip.bottom
            }

            font {
                family: "Microsoft YaHei"
                pixelSize: 18
            }

            echoMode: TextInput.Password
        }

        Text {
            id: regemilTip
            text: qsTr("请输入邮箱")
            visible: !mainWindow.tabIsLogin
            color: "#1A1A1A"
            opacity: 0.3  // 添加透明度，0-1之间
            font {
                family: "Microsoft YaHei"
                pixelSize: 14
            }
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: regPwdInBox.bottom
                topMargin: 20
            }
        }

        Antui.FramelessLineEdit {
            id: regemilInBox
            visible: !mainWindow.tabIsLogin
            width: parent.width * 0.56
            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                top: regemilTip.bottom
            }

            font {
                family: "Microsoft YaHei"
                pixelSize: 18
            }
        }

        Button {
            id: regButton
            visible: !mainWindow.tabIsLogin
            width: parent.width * 0.56
            height: 45
            text: "提交"

            anchors {
                left: parent.left
                leftMargin: parent.width*0.22
                bottom: parent.bottom
                bottomMargin: parent.height * 0.14
            }

            background: Rectangle {
                radius: 10  // 圆角
                color: regButton.down ? "#3d8b40" :  "#5CBF60"

                border.width: 0
            }

            contentItem: Text {
                text: regButton.text
                color: "white"  // 白色文字
                font {
                    pixelSize: 20
                    family: "Microsoft YaHei"
                    bold: true
                }
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
            }

            onClicked: {
                console.log("登录按钮被点击")
            }
        }
    }
}
