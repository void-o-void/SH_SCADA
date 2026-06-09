cmake_minimum_required(VERSION 3.16)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTORCC ON)

project(HMI LANGUAGES CXX)

find_package(Qt6 6.5 REQUIRED COMPONENTS Quick QuickControls2)
qt_standard_project_setup(REQUIRES 6.5)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)

qt_add_executable(hmi_app
    src/main.cpp
    src/devicelistmodel.cpp
    src/thememgr.cpp
    src/usermgr.hpp
)

include_directories( hmi_app
    src
    src/common
)

# 添加 QML 模块
qt_add_qml_module(hmi_app
    URI HMI
    VERSION 1.0
    OUTPUT_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}/HMI
    QML_FILES
        src/qml/Login.qml
        src/qml/MainWindow.qml
        src/qml/SidebarMenu.qml
        src/qml/Navigation.qml
        src/qml/TopFrame.qml
        antui/FramelessCombox.qml
        antui/FramelessLineEdit.qml
        antui/IconButton.qml
    RESOURCES
        resource/resource.qrc
)

# 链接库
target_link_libraries(hmi_app PRIVATE
    Qt6::Core
    Qt6::Quick
    Qt6::Qml
    Qt6::QuickControls2
)

# 设置可执行文件名
set_target_properties(hmi_app PROPERTIES
    OUTPUT_NAME "hmi"
)

