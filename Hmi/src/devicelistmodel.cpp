#include "devicelistmodel.h"
#include <QRandomGenerator>
#include <QDebug>

CDeviceListModel::CDeviceListModel(QObject *parent)
    : QAbstractListModel(parent)
{
}

int CDeviceListModel::rowCount(const QModelIndex &parent) const
{
    Q_UNUSED(parent)
    return m_devices.size();
}

QVariant CDeviceListModel::data(const QModelIndex &index, int role) const
{
    if (!index.isValid() || index.row() >= m_devices.size())
        return QVariant();

    const SDeviceItem &item = m_devices.at(index.row());

    switch (role) {
    case Qt::DisplayRole:
        if(item.level == 1) {
            return item.group;
        }
        return item.name;
    case Qt::EditRole:
        return item.name;
    case IconRole:
        return item.iconbuf;
    case LevelRole:
        return item.level;
    case ParentNameRole:
        return item.group;
    case ChildsNameRole:
        return item.product;
    default:
        return QVariant();
    }
}

QHash<int, QByteArray> CDeviceListModel::roleNames() const
{
    QHash<int, QByteArray> roles;
    roles[Qt::DisplayRole] = "display";
    roles[IconRole] = "icon";
    roles[LevelRole] = "level";
    roles[ParentNameRole] = "parentName";
    roles[ChildsNameRole] = "childsName";
    return roles;
}

void CDeviceListModel::syncUpdate()
{
    beginResetModel();

    // 清空现有数据
    m_devices.clear();

    // 定义可能的设备分组
    QStringList groups = {
        "会议室设备",
        "办公设备",
        "服务器设备",
        "网络设备",
        "安防设备"
    };

    // 定义可能的设备类型
    QStringList deviceTypes = {
        "投影仪", "电视", "音响", "摄像头", "路由器",
        "交换机", "服务器", "打印机", "扫描仪", "电脑"
    };

    // 定义可能的品牌
    QStringList brands = {
        "华为", "海康威视", "大华", "小米", "TP-LINK",
        "H3C", "戴尔", "联想", "惠普", "苹果"
    };

    // 用于存储设备数据
    QVector<SDeviceItem> allDevices;

    // 生成100个设备
    for (int i = 1; i <= 100; ++i) {
        SDeviceItem device;
        device.level = 2;  // 真实设备是2级

        // 随机分配分组
        int groupIndex = QRandomGenerator::global()->bounded(groups.size());
        device.group = groups[groupIndex];

        // 随机选择设备类型
        int typeIndex = QRandomGenerator::global()->bounded(deviceTypes.size());
        QString deviceType = deviceTypes[typeIndex];

        // 随机选择品牌
        int brandIndex = QRandomGenerator::global()->bounded(brands.size());
        QString brand = brands[brandIndex];

        // 生成设备名称
        device.name = QString("%1-%2-%3").arg(brand).arg(deviceType).arg(i);

        // 生成设备描述
        device.descr = QString("这是%1%2，型号%3").arg(brand).arg(deviceType).arg(i);

        // 设备图标（模拟）
        device.iconbuf = QString("device_icon_%1").arg(deviceType);

        // 产品信息
        device.product = QString("Product-%1").arg(1000 + i);

        allDevices.append(device);
    }

    // 统计每个分组的设备数量
    QHash<QString, int> groupDeviceCounts;
    for (const auto &device : allDevices) {
        groupDeviceCounts[device.group]++;
    }

    // 按分组1-分组下的设备1，分组下的设备2...分组2...的顺序存储
    for (const auto &group : groups) {
        // 如果这个分组有设备，先添加分组
        if (groupDeviceCounts.contains(group)) {
            SDeviceItem groupItem;
            groupItem.level = 1;  // 分组是1级
            groupItem.name = group;
            groupItem.group = group;
            groupItem.isExpeand = true;  // 默认展开
            groupItem.descr = QString("%1（共%2个设备）").arg(group).arg(groupDeviceCounts[group]);
            m_devices.append(groupItem);

            // 添加这个分组下的所有设备
            for (const auto &device : allDevices) {
                if (device.group == group) {
                    m_devices.append(device);
                }
            }
        }
    }

    endResetModel();

    qDebug() << "同步更新完成，共生成" << m_devices.size() << "个项目";
    qDebug() << "其中1级菜单（分组）:" << groupDeviceCounts.size() << "个";
    qDebug() << "2级菜单（设备）: 100个";
}
