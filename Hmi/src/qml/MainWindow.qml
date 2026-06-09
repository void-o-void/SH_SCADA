import QtQuick
import QtQuick.Controls
import QtQuick.Window

ApplicationWindow {
    id: window
    width: 1920
    height: 1080
    visible: true
    color: "#FFFFFF"
    //flags: Qt.Window | Qt.FramelessWindowHint

    TopFrame {
        id:topframe
    }

    Navigation {
        anchors.top: topframe.bottom
    }
}
