// SidebarMenu.qaml
import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import QtQuick.Shapes
import HMI.Device 1.0  // 导入你的模型模块

Rectangle {
    id: root
    width: 200
    height: parent.height
    color: "#001529"  // Ant Design 深色主题

    property DeviceListModel deviceListModel: null  // 用于接收C++模型的属性
    property var expandedGroups: ({})  // 存储分组展开状态的字典
    property string selectedMenuItem: ""  // 当前选中的菜单项key

    signal menuItemClicked(string key, string label)
    signal subMenuToggled(string key, bool expanded)

    // 初始化分组展开状态
    function initExpandedGroups() {
        if (!deviceListModel) return

        for (var i = 0; i < deviceListModel.rowCount(); i++) {
            var itemData = deviceListModel.data(deviceListModel.index(i, 0), deviceListModel.LevelRole)
            if (itemData === 1) {  // 一级菜单
                var groupName = deviceListModel.data(deviceListModel.index(i, 0), Qt.DisplayRole)
                expandedGroups[groupName] = true  // 默认展开
            }
        }
    }

    // 检查分组是否展开
    function isGroupExpanded(groupName) {
        return expandedGroups[groupName] === true
    }

    // 处理分组点击
    function toggleGroup(index) {
        var groupName = deviceListModel.data(deviceListModel.index(index, 0), Qt.DisplayRole)
        if (groupName) {
            var currentState = expandedGroups[groupName] || false
            expandedGroups[groupName] = !currentState
            subMenuToggled(groupName, expandedGroups[groupName])

            // 选中这个分组
            selectedMenuItem = groupName
            root.menuItemClicked(groupName, groupName)

            // 强制刷新ListView
            var oldModel = menuListView.model
            menuListView.model = null
            menuListView.model = oldModel
        }
    }

    onDeviceListModelChanged: {
        if (deviceListModel) {
            // 监听模型数据变化
            deviceListModel.rowsInserted.connect(initExpandedGroups)
            deviceListModel.modelReset.connect(initExpandedGroups)

            // 初始化展开状态
            initExpandedGroups()
        }
    }

    // 顶部搜索区域
    Rectangle {
        id: toggleArea
        width: parent.width
        height: 48
        color: "#011520"

        Row {
            anchors.centerIn: parent
            spacing: 10

            Button {
                text: "同步设备"
                onClicked: {
                    if (deviceListModel) {
                        deviceListModel.syncUpdate()
                    }
                }
            }
        }
    }

    // 菜单项容器
    ListView {
        id: menuListView
        width: parent.width
        anchors.top: toggleArea.bottom
        anchors.bottom: parent.bottom
        anchors.topMargin: 8
        clip: true
        spacing: 0

        model: deviceListModel  // 使用C++模型

        delegate: Item {
            id: delegateItem
            width: menuListView.width

            // 控制是否显示二级设备
            readonly property bool isGroupItem: level === 1
            readonly property string currentGroup: isGroupItem ? display : parentName
            readonly property bool isExpanded: root.isGroupExpanded(currentGroup)
            readonly property bool shouldShow: isGroupItem || isExpanded
            readonly property bool isSelected: root.selectedMenuItem === (isGroupItem ? true : false)

            // 一级菜单高度50，二级设备高度40
            height: shouldShow ? (isGroupItem ? 50 : 40) : 0
            visible: shouldShow

            Rectangle {
                id: delegateBackground
                anchors.fill: parent
                color: {
                    if (isSelected) {
                        return "#1890ff"  // 选中状态的背景色
                    } else if (level === 1) {
                        return mouseArea.containsMouse ? "#00284d" : "transparent"
                    } else {
                        return mouseArea.containsMouse ? "#1890ff" : "transparent"
                    }
                }

                // 左侧占位，用于缩进
                Item {
                    id: leftSpacer
                    width: level === 1 ? 16 : 40
                    height: parent.height
                }

                // 二级菜单图标
                Rectangle {
                    id: iconContainer
                    width: 20
                    height: 20
                    radius: 3
                    color: level === 2 ? (isSelected ? "#ffffff" : "#1890ff") : "transparent"
                    visible: level === 2
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.left: leftSpacer.right

                    Text {
                        anchors.centerIn: parent
                        text: {
                            // 使用图标库中的字符图标
                            if (icon && typeof icon === "string") {
                                if (icon.includes("desktop")) return "💻"
                                else if (icon.includes("server")) return "🖥️"
                                else if (icon.includes("network")) return "🌐"
                                else if (icon.includes("printer")) return "🖨️"
                                else if (icon.includes("camera")) return "📷"
                                else return "📱"
                            }
                            return "📱"
                        }
                        color: (level === 2 && isSelected) ? "#1890ff" : "white"
                        font.pixelSize: 12
                    }
                }

                // 菜单文本
                Text {
                    id: menuText
                    text: display
                    color: {
                        if (isSelected) {
                            return "#ffffff"
                        } else if (level === 1) {
                            return mouseArea.containsMouse ? "#ffffff" : "#b2b2b2"
                        } else {
                            return mouseArea.containsMouse ? "#ffffff" : "#d9d9d9"
                        }
                    }
                    font.pixelSize: level === 1 ? 16 : 14
                    font.bold: level === 1
                    elide: Text.ElideRight
                    anchors.verticalCenter: parent.verticalCenter

                    // 如果是二级菜单，图标在文本左边
                    anchors.left: level === 2 ? iconContainer.right : leftSpacer.right
                    anchors.leftMargin: level === 2 ? 8 : 0
                }

                // 一级菜单折叠按钮（只在右侧显示）
                Item {
                    id: expandCollapseButton
                    width: 20
                    height: 20
                    visible: level === 1
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.right: parent.right
                    anchors.rightMargin: 8

                    Text {
                        anchors.centerIn: parent
                        text: isExpanded ? "▼" : "▶"
                        color: isSelected ? "#ffffff" : "#b2b2b2"
                        font.pixelSize: 10
                    }
                }

                // 左侧选中指示条
                Rectangle {
                    width: 3
                    height: parent.height
                    color: "#1890ff"
                    visible: isSelected
                }

                // 分隔线
                Rectangle {
                    anchors.bottom: parent.bottom
                    width: parent.width
                    height: 1
                    color: "#003261"
                    visible: level === 2 || isExpanded
                }
            }

            MouseArea {
                id: mouseArea
                anchors.fill: parent
                hoverEnabled: true
                //cursorShape: Qt.PointingHandCursor
                onClicked: {
                    if (level === 1) {
                        // 一级菜单点击 - 切换展开/折叠状态
                        toggleGroup(index)
                    } else {
                        // 二级设备点击
                        var deviceName = deviceListModel.data(deviceListModel.index(index, 0), Qt.DisplayRole)
                        var parentGroup = deviceListModel.data(deviceListModel.index(index, 0), deviceListModel.ParentNameRole)

                        // 选中这个设备
                        root.selectedMenuItem = deviceName
                        root.menuItemClicked(parentGroup, deviceName)
                    }
                }
            }
        }
    }

    // 宽度动画
    Behavior on width {
        NumberAnimation {
            duration: 200
            easing.type: Easing.InOutQuad
        }
    }
}
