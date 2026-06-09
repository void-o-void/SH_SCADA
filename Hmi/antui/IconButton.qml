import QtQuick
import QtQuick.Controls
import HMI.Core 1.0

Button {
    id: root

    property string buttonName: "home"            // 按钮名称
    property string descr: "1"                     // 功能描述

    property string iconfamily: "regular"         // 所属字体库
    property string iconCode: "1"                     // 图标码
    property int iconSize: 16                     // 图标大小

    property string textText: "按钮"                   // 文本
    property int textSize: 14                     // 文字大小
    property string textFamily: "1"                   // 文本字体

    property bool isSelect: false                 // 是否选中(只有多个组合，需要保持状态，例如导航菜单才会使用此属性)

    property int leftSpacing: 8                   // 左边距
    property int rightSpacing: 8                  // 右边距
    property bool iconLeft: true                  // 图标在左侧（true）还是右侧（false）

    property int radius: 0                        // 圆角

    // 内容区域
    contentItem: Row {
        id: contentRow
        spacing: 4
        layoutDirection: root.iconLeft ? Qt.LeftToRight : Qt.RightToLeft

        // 图标
        Text {
            id: iconText
            text: root.iconCode
            font.family: root.iconfamily
            font.pixelSize: root.iconSize

            // 图标颜色绑定：根据按钮状态使用不同的主题颜色
            color: {
                if (root.down)
                    return CThemeMgr.buttonSelectedIconColor
                else if (root.hovered)
                    return CThemeMgr.buttonHoverIconColor
                else if (root.isSelect)
                    return CThemeMgr.buttonSelectedIconColor
                else
                    return CThemeMgr.buttonIconColor
            }

            verticalAlignment: Text.AlignVCenter
            renderType: Text.NativeRendering
            visible: text !== ""
        }

        // 文本
        Text {
            id: labelText
            text: root.textText
            font.pixelSize: root.textSize
            font.family: root.textFamily

            // 文本颜色绑定：根据按钮状态使用不同的主题颜色
            color: {
                if (root.down)
                    return CThemeMgr.buttonSelectedTextColor
                else if (root.hovered)
                    return CThemeMgr.buttonHoverTextColor
                else if (root.isSelect)
                    return CThemeMgr.buttonSelectedTextColor
                else
                    return CThemeMgr.buttonTextColor
            }

            verticalAlignment: Text.AlignVCenter
            horizontalAlignment: Text.AlignHCenter
        }
    }

    // 按钮背景样式
    background: Rectangle {
        implicitWidth: 100
        implicitHeight: 40

        // 背景颜色绑定：根据按钮状态使用不同的主题颜色
        color: {
            if (root.down)
                return CThemeMgr.buttonSelectedBgColor
            else if (root.hovered)
                return CThemeMgr.buttonHoverBgColor
            else if (root.isSelect)
                return CThemeMgr.buttonSelectedBgColor
            else
                return CThemeMgr.buttonBgColor
        }

        border.color: root.activeFocus ? CThemeMgr.selectedBgColor : CThemeMgr.borderColor
        border.width: root.activeFocus ? 2 : 1
        radius: root.radius
    }

    // 添加边距控制
    leftPadding: root.leftSpacing
    rightPadding: root.rightSpacing
}
