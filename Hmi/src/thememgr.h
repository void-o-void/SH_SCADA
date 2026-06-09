#ifndef THEMEMGR_H
#define THEMEMGR_H

#include <QObject>
#include <QFontDatabase>
#include <QDebug>
#include "common/common.h"

class CThemeMgr : public QObject
{
    Q_OBJECT
    DECLARE_SINGLETON(CThemeMgr)

    DECLARE_QML_QPROPERTY(bgColor, QString)
    DECLARE_QML_QPROPERTY(selectedBgColor, QString)
    DECLARE_QML_QPROPERTY(fgColor, QString)
    DECLARE_QML_QPROPERTY(selectedFgColor, QString)
    DECLARE_QML_QPROPERTY(borderColor, QString)
    DECLARE_QML_QPROPERTY(topBorderColor, QString)
    DECLARE_QML_QPROPERTY(leftBorderColor, QString)
    DECLARE_QML_QPROPERTY(rightBorderColor, QString)
    DECLARE_QML_QPROPERTY(bottmBorderColor, QString)


    // Button 专属定制属性
    DECLARE_QML_QPROPERTY(buttonBgColor, QString)           // 1、背景颜色
    DECLARE_QML_QPROPERTY(buttonSelectedBgColor, QString)   // 2、选中时的背景颜色
    DECLARE_QML_QPROPERTY(buttonHoverBgColor, QString)      // 3、悬浮时的背景颜色
    DECLARE_QML_QPROPERTY(buttonTextColor, QString)         // 4、默认文本颜色
    DECLARE_QML_QPROPERTY(buttonHoverTextColor, QString)    // 5、悬浮文本颜色
    DECLARE_QML_QPROPERTY(buttonSelectedTextColor, QString) // 6、选中时的文本颜色
    DECLARE_QML_QPROPERTY(buttonIconColor, QString)         // 7、默认图标颜色
    DECLARE_QML_QPROPERTY(buttonHoverIconColor, QString)    // 8、悬浮时图标颜色
    DECLARE_QML_QPROPERTY(buttonSelectedIconColor, QString) // 9、选中时候的图标颜色

    //字体管理
    FONT_PROPERTY(faBrandsRegular);
    FONT_PROPERTY(faDuotoneLight);
    FONT_PROPERTY(faDuotoneRegular);
    FONT_PROPERTY(faDuotoneSolid);
    FONT_PROPERTY(faDuotoneThin);
    FONT_PROPERTY(faProLight);
    FONT_PROPERTY(faProRegular);
    FONT_PROPERTY(faProSolid);
    FONT_PROPERTY(faProThin);
    FONT_PROPERTY(faSharpLight);
    FONT_PROPERTY(faSharpRegular);
    FONT_PROPERTY(faSharpSolid);
    FONT_PROPERTY(faSharpThin);
    FONT_PROPERTY(faSharpDuotoneLight);
    FONT_PROPERTY(faSharpDuotoneRegular);
    FONT_PROPERTY(faSharpDuotoneSolid);
    FONT_PROPERTY(faSharpDuotoneThin);


    enum ETheme {
        Light = 0x00,
        Drak
    };
    Q_ENUM(ETheme)
public:
    Q_INVOKABLE ETheme currentTheme() { return m_currentTheme; }
    Q_INVOKABLE void setTheme(ETheme theme);

protected:
    void setButtonTheme(ETheme theme);
    void loadFonts();

signals:

private:
    explicit CThemeMgr(QObject *parent = nullptr);
    enum ETheme m_currentTheme = Light;
};

#endif // THEMEMGR_H
