import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

TextField {
    id: lineEdit
    width: 200
    height: 40

    // 设置边框
    background: Rectangle {
        color: "transparent"
        border.color: "#A6A6A6"
        border.width: 0
        Rectangle {
            anchors.bottom: parent.bottom
            height: 0.8
            width: lineEdit.width
            color: "#1A1A1A"
            opacity: 0.3  // 添加透明度，0-1之间
        }
    }

    // 设置文本对齐
    horizontalAlignment: Text.AlignLeft
    verticalAlignment: Text.AlignBottom

    // 设置内边距
    topPadding: 0
    bottomPadding: 4
    leftPadding: 0
    rightPadding: 0
}
