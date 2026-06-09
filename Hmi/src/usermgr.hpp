#ifndef USERMGR_H
#define USERMGR_H

#include <QObject>
#include <QJsonDocument>
#include <QJsonObject>
#include <QDebug>
#include "common/common.h"
#include "common/httplib.h"

class CUserMgr : public QObject
{
    Q_OBJECT
    DECLARE_SINGLETON(CUserMgr)
public:
    Q_PROPERTY(QString token READ token NOTIFY tokenChanged)
    Q_PROPERTY(bool loggedIn READ loggedIn NOTIFY loggedInChanged)

    Q_INVOKABLE bool login(const QString& user, const QString& pwd) {
        QJsonObject body;
        body["username"] = user;
        body["password"] = pwd;
        std::string jsonBody = QJsonDocument(body).toJson(QJsonDocument::Compact).toStdString();

        qDebug() << "[Login]" << user << "@" << m_host;

        httplib::Client cli(m_host.toStdString(), 8080);
        cli.set_connection_timeout(3);

        auto res = cli.Post("/api/login",{{"Content-Type", "application/json"}},jsonBody, "application/json");

        if (!res) {
            qWarning() << "[Login] 连接失败";
            return false;
        }

        qDebug() << "[Login] HTTP" << res->status;

        if (res->status != 200) {
            qWarning() << "[Login] 失败:" << res->body.c_str();
            return false;
        }

        QJsonDocument doc = QJsonDocument::fromJson(QByteArray::fromStdString(res->body));
        m_token = doc.object()["token"].toString();
        m_loggedIn = true;

        qDebug() << "[Login] 成功";
        emit tokenChanged();
        emit loggedInChanged();
        return true;
    }

    void setHost(const QString& host) { m_host = host; }
    QString token() const { return m_token; }
    bool loggedIn() const { return m_loggedIn; }

signals:
    void tokenChanged();
    void loggedInChanged();

private:
    explicit CUserMgr(QObject *parent = nullptr) : QObject(parent) {}
    QString m_host = "127.0.0.1";
    QString m_token;
    bool m_loggedIn = false;
};

#endif // USERMGR_H
