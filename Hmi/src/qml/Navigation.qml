import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import HMI.Core 1.0
import "../../antui" as Antui

Rectangle {
    id: root
    width: 240
    height: 1020
    visible: true
    color: "#04142C"

    Column {
        anchors.fill: parent
        spacing: 8
        padding: 16

        // 第一个按钮 - 首页
        Antui.IconButton {
            id: homeButton
            width: parent.width - 32
            height: 50
            buttonName: "home"
            descr: "返回首页"
            iconCode: "\uf17b"
            iconfamily: CThemeMgr.faBrandsRegular
            iconSize: 20
            textText: "首页"
            textSize: 16
            leftSpacing: 16
            rightSpacing: 16
            radius: 8
            isSelect: true  // 默认选中首页
        }

        // 第二个按钮 - 设置
        Antui.IconButton {
            id: settingsButton
            width: parent.width - 32
            height: 50
            buttonName: "settings"
            descr: "系统设置"
            iconCode: "\uf013"
            iconfamily: CThemeMgr.faBrandsRegular
            iconSize: 20
            textText: "设置"
            textSize: 16
            leftSpacing: 16
            rightSpacing: 16
            radius: 8
        }

        // 第三个按钮 - 搜索
        Antui.IconButton {
            id: searchButton
            width: parent.width - 32
            height: 50
            buttonName: "search"
            descr: "搜索功能"
            iconCode: "\uf002"
            iconfamily: CThemeMgr.faBrandsRegular
            iconSize: 20
            textText: "搜索"
            textSize: 16
            leftSpacing: 16
            rightSpacing: 16
            radius: 8
        }

        // 第四个按钮 - 消息
        Antui.IconButton {
            id: messageButton
            width: parent.width - 32
            height: 50
            buttonName: "message"
            descr: "查看消息"
            iconCode: "\uf0e0"
            iconfamily: CThemeMgr.faBrandsRegular
            iconSize: 20
            textText: "消息"
            textSize: 16
            leftSpacing: 16
            rightSpacing: 16
            radius: 8
        }

        // 第五个按钮 - 帮助
        Antui.IconButton {
            id: helpButton
            width: parent.width - 32
            height: 50
            buttonName: "help"
            descr: "帮助中心"
            iconCode: "\uf059"
            iconfamily: CThemeMgr.faBrandsRegular
            iconSize: 20
            textText: "帮助"
            textSize: 16
            leftSpacing: 16
            rightSpacing: 16
            radius: 8
        }

        // 添加一个占位项，让按钮在顶部显示
        Item {
            Layout.fillHeight: true
        }
    }
}
