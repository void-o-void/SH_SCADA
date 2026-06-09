import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

// 自定义下拉框组件（支持手动输入）
ComboBox {
    id: customComboBox
    width: 200
    height: 40

    // 添加可编辑属性
    editable: true

    // 隐藏默认边框，仅保留底部边框
    background: Rectangle {
        color: "transparent"
        border.color: "#A6A6A6"
        border.width: 0
        Rectangle {
            anchors.bottom: parent.bottom
            height: 1
            width: customComboBox.width
            color: "#1A1A1A"
            opacity: 0.3  // 添加透明度，0-1之间
        }
    }

    indicator: Item {
        id: arrowIcon
        width: 25
        height: 25
        anchors.right: parent.right
        anchors.rightMargin: 5
        anchors.bottom: parent.bottom
        anchors.bottomMargin: 5

        // 可以添加旋转动画
        rotation: customComboBox.popup.opened ? 180 : 0

        Behavior on rotation {
            NumberAnimation { duration: 200 }
        }

        Image {
            anchors.fill: parent
            source: "/icon/down.svg"
            fillMode: Image.PreserveAspectFit
            sourceSize.width: parent.width
            sourceSize.height: parent.height
        }
    }

    // ========== 显示当前选中项的样式 ==========
    contentItem: TextField {
        text: customComboBox.editable ? customComboBox.editText : customComboBox.displayText
        font: customComboBox.font
        color: "#383838"
        verticalAlignment: Text.AlignBottom
        horizontalAlignment: Text.AlignLeft

        // 关键设置：左边距为0，对齐下边框左边
        leftPadding: 0
        rightPadding: customComboBox.indicator.width + 5

        // 确保文本框背景透明
        background: Rectangle {
            color: "transparent"
            border.width: 0
        }

        // 可编辑时的额外配置
        readOnly: !customComboBox.editable
        selectByMouse: customComboBox.editable

        onAccepted: {
            if (customComboBox.editable) {
                customComboBox.editText = text
                customComboBox.currentIndex = -1
            }
        }
    }

    // ========== 下拉列表的样式美化 ==========
    popup: Popup {
        y: customComboBox.height + 2
        width: customComboBox.width
        implicitHeight: Math.min(contentItem.implicitHeight, 200)
        padding: 0

        background: Rectangle {
            color: "#FFFFFF"
            border.color: "#EEEEEE"
            border.width: 0
            radius: 10
        }

        contentItem: ListView {
            clip: true
            implicitHeight: contentHeight
            model: customComboBox.popup.visible ? customComboBox.delegateModel : null
            currentIndex: customComboBox.highlightedIndex

            ScrollIndicator.vertical: ScrollIndicator { }

            highlight: Rectangle {
                color: "#F5F5F5"
                radius: 2
            }
        }
    }

    // ========== 下拉选项的样式 ==========
    delegate: ItemDelegate {
        width: customComboBox.width
        height: 36

        // 关键设置：选项文本也左对齐，无左边距
        leftPadding: 0

        contentItem: Text {
            text: modelData
            font.pixelSize: 14
            color: "#333333"
            verticalAlignment: Text.AlignVCenter
            elide: Text.ElideRight
        }

        background: Rectangle {
            color: customComboBox.highlightedIndex === index ? "#F5F5F5" : "transparent"
        }
    }
}
