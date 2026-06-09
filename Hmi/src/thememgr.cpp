#include "thememgr.h"

CThemeMgr::CThemeMgr(QObject *parent)
    : QObject{parent}
{
    setTheme(Light);
    loadFonts();
}

void CThemeMgr::setTheme(ETheme theme) {
    switch (theme) {
    case Drak:
        set_bgColor("#04142C");
        set_selectedBgColor("#1677FF");
        set_fgColor("#A6ADB4");
        set_selectedFgColor("#FFFFFF");
        break;
    case Light:
        set_bgColor("#FFFFFF");
        set_selectedBgColor("#E6F4FF");
        set_fgColor("#1F1F1F");
        set_selectedFgColor("#1677FF");
        break;
    default:
        return;
    }
    m_currentTheme = theme;
    setButtonTheme(theme);
}

void CThemeMgr::setButtonTheme(ETheme theme) {
    switch (theme) {
    case Drak:
        set_buttonBgColor("#04142C");
        set_buttonSelectedBgColor("#1677FF");
        set_buttonHoverBgColor("#1F1F1F");
        set_buttonTextColor("#A6ADB4");
        set_buttonHoverTextColor("#FFFFFF");
        set_buttonSelectedTextColor("#FFFFFF");
        set_buttonIconColor("#A6ADB4");
        set_buttonHoverIconColor("#FFFFFF");
        set_buttonSelectedIconColor("#FFFFFF");
        break;
    case Light:
        set_buttonBgColor("#FFFFFF");
        set_buttonSelectedBgColor("#E6F4FF");
        set_buttonHoverBgColor("#F5F5F5");
        set_buttonTextColor("#1F1F1F");
        set_buttonHoverTextColor("#1F1F1F");
        set_buttonSelectedTextColor("#1677FF");
        set_buttonIconColor("#595959");
        set_buttonHoverIconColor("#1F1F1F");
        set_buttonSelectedIconColor("#1677FF");
        break;
    default:
        break;
    }
}

void CThemeMgr::loadFonts()
{
    // 定义字体文件路径数组
    struct FontInfo {
        QString* var;
        QString path;
    };

    // 使用图片中的实际文件名（下划线格式，无空格）
    FontInfo fonts[] = {
        {&m_faBrandsRegular, ":/fontawesome6.72pro/Brands_Regular_400.otf"},
        {&m_faDuotoneLight, ":/fontawesome6.72pro/Duotone_Light_300.otf"},
        {&m_faDuotoneRegular, ":/fontawesome6.72pro/Duotone_Regular_400.otf"},
        {&m_faDuotoneSolid, ":/fontawesome6.72pro/Duotone_Solid_900.otf"},
        {&m_faDuotoneThin, ":/fontawesome6.72pro/Duotone_Thin_100.otf"},
        {&m_faProLight, ":/fontawesome6.72pro/Pro_Light_300.otf"},
        {&m_faProRegular, ":/fontawesome6.72pro/Pro_Regular_400.otf"},
        {&m_faProSolid, ":/fontawesome6.72pro/Pro_Solid_900.otf"},
        {&m_faProThin, ":/fontawesome6.72pro/Pro_Thin_100.otf"},
        {&m_faSharpLight, ":/fontawesome6.72pro/Sharp_Light_300.otf"},
        {&m_faSharpRegular, ":/fontawesome6.72pro/Sharp_Regular_400.otf"},
        {&m_faSharpSolid, ":/fontawesome6.72pro/Sharp_Solid_900.otf"},
        {&m_faSharpThin, ":/fontawesome6.72pro/Sharp_Thin_100.otf"},
        {&m_faSharpDuotoneLight, ":/fontawesome6.72pro/SharpDuotone_Light_300.otf"},
        {&m_faSharpDuotoneRegular, ":/fontawesome6.72pro/SharpDuotone_Regular_400.otf"},
        {&m_faSharpDuotoneSolid, ":/fontawesome6.72pro/SharpDuotone_Solid_900.otf"},
        {&m_faSharpDuotoneThin, ":/fontawesome6.72pro/SharpDuotone_Thin_100.otf"}
    };

    // 批量加载字体
    for (const auto& font : fonts) {
        int fontId = QFontDatabase::addApplicationFont(font.path);
        if(fontId < 0) {
            qDebug() << "Failed to load font:" << font.path;
            continue;
        }
        *font.var = QFontDatabase::applicationFontFamilies(fontId).at(0);
        qDebug() << "Successfully loaded font:" << *font.var << "  path:" << font.path;
    }
}

