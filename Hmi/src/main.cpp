#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQuickStyle>
#include <QtGui>
#include <QIcon>

#include "devicelistmodel.h"
#include "src/usermgr.hpp"
#include "src/thememgr.h"

int main(int argc, char *argv[])
{
    QGuiApplication app(argc, argv);

    qmlRegisterSingletonInstance<CUserMgr>("HMI.Core", 1, 0, "UserMgr", CUserMgr::instance());
    qmlRegisterSingletonInstance<CThemeMgr>("HMI.Core", 1, 0, "CThemeMgr", CThemeMgr::instance());
    qmlRegisterType<CDeviceListModel>("HMI.Device", 1, 0, "DeviceListModel");

    QQuickStyle::setStyle("Basic");
    app.setWindowIcon(QIcon("qrc:/icon/iot.svg"));

    QQmlApplicationEngine engine;
    QObject::connect(
        &engine,
        &QQmlApplicationEngine::objectCreationFailed,
        &app,
        []() { QCoreApplication::exit(-1); },
        Qt::QueuedConnection);
    engine.loadFromModule("HMI", "Login");
    return app.exec();
}
