#ifndef COMMON_H
#define COMMON_H

#define DECLARE_SINGLETON(ClassName) \
public: \
    static ClassName* instance(){ \
        static ClassName instance; \
        return &instance; \
    }; \
private: \
    ClassName(const ClassName&) = delete; \
    ClassName& operator=(const ClassName&) = delete; \
    ClassName(ClassName&&) = delete; \
    ClassName& operator=(ClassName&&) = delete;


#define DECLARE_READ_PROPERTY(name, type) \
private: \
    type m_##name; \
    Q_PROPERTY(type m_##name READ get_##name CONSTANT)\
public: \
    inline type get_##name() const { return m_##name;}

//字体宏
#define FONT_PROPERTY(name) \
Q_PROPERTY(QString name READ name CONSTANT) \
    QString name() const { return m_##name; } \
    QString m_##name


#define DECLARE_QML_QPROPERTY(name, type) \
public: \
    type m_##name;\
    Q_PROPERTY(type name READ get_##name WRITE set_##name NOTIFY name##Changed) \
    Q_SIGNAL void name##Changed();\
    Q_SLOT void set_##name(const type& value) { if(m_##name != value){ m_##name = value; emit name##Changed();} } \
    type get_##name() const { return m_##name; }

#endif // COMMON_H
