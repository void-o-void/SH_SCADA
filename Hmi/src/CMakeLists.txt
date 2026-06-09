set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTORCC ON)

find_package(Qt6 6.5 REQUIRED COMPONENTS Quick QuickControls2)
qt_standard_project_setup(REQUIRES 6.5)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)

qt_add_executable(hmi_app
    main.cpp
    devicelistmodel.cpp
    usermgr.hpp
)

# 添加 QML 模块
qt_add_qml_module(hmi_app
    URI HMI
    VERSION 1.0
    OUTPUT_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}/HMI
    QML_FILES
        qml/Login.qml
        qml/Main.qml
        qml/SidebarMenu.qml
        ../resource/resource.qrc
)

# 链接库
target_link_libraries(hmi_app PRIVATE
    Qt6::Core
    Qt6::Quick
    Qt6::Qml
    Qt6::QuickControls2
    antui
)

# 设置可执行文件名
set_target_properties(hmi_app PROPERTIES
    OUTPUT_NAME "hmi"
)
