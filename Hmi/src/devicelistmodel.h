#ifndef CDeviceListModel_H
#define CDeviceListModel_H

#include <QObject>
#include <QAbstractListModel>
#include <QVector>
#include <QHash>

class CDeviceListModel : public QAbstractListModel
{
    Q_OBJECT
public:
    struct SDeviceInfo {
        QString name;  //用于显示的文本，如果是1级菜单此处为分组名字
        QString descr;
        QString iconbuf;
        QString group;
        QString product;
    };

    struct SDeviceItem : public SDeviceInfo {
        int level;    //菜单级别，如果是真实设备则属于2级设备，如果是生成的虚拟分组则是1级
        bool isExpeand; //只对1级菜单有效，标识是否展开状态
    };

    enum CDeviceListModelRole {
        IconRole = Qt::UserRole + 1,
        LevelRole,
        ParentNameRole,
        ChildsNameRole
    };
    Q_ENUM(CDeviceListModelRole)

    explicit CDeviceListModel(QObject *parent = nullptr);

    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QHash<int, QByteArray> roleNames() const override;

    Q_INVOKABLE void syncUpdate();

signals:

private:
    QVector<SDeviceItem> m_devices;
};

#endif // CDeviceListModel_H
