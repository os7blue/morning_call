###### morning_call 之前比较火的每天通过微信公众号给女朋友推送天气，纪念日，一句话的golang实现

效果截图：

![301664737018_ pic](https://user-images.githubusercontent.com/33112372/193471271-d707397f-abc4-497d-bc3f-924afbc9243b.jpg)


为啥要写这玩意儿：
        
    之前这玩意儿比较火，我也给我家领导整了一个。
    但是吧，之前用的别人写的。
    他的每日一句话接口用的网易家的，毕竟"网抑云"嘛。
    于是呼，我遭到了背刺。
    网易家的那个接口，每日一句话就很灵性。
    比如：
        在假期刚结束的时候，推送的是：没有人比过完假期的人更需要假期了。
        在我家领导最忙的时候，推送的是：我是一个没有感情的工作机器。
        在我家领导差点因为疫情回不了家的时候，推送的是：你买到火车票了吗？
    
    当然，我晓得，这些事儿也不能怪接口。
    但是！
    为了我的生命安全，我决定自己写实现并把一句话接口换成一言。
    所以这个仓库诞生了。

使用步骤：
    
    1、fork这个仓库，并设置私有。
    2、搞一个服务号（测试号怎么搞不用我教了吧），并建立对应的消息模板。
    3、去和风天气注册一个开发者（这个天气平台api挺全的）。
    4、打开config.ini填写相应的配置（ini文件注释写的很全，别问。）
    5、使用github action构建每日任务。


配置文件内容也贴一下：
    
    [Wechat]
    #公众号appid
    app-id =
    #公众号key
    app-secret =
    #消息模板id
    template-id =
    #接受信息用户openId,数组，以,分割
    user =
    
    
    [Weather]
    #和风天气 key
    key =
    #和风天气定位，可数组，以,分隔。 如 北京市,上海市
    region =
    
    [Hitokoto]
    #句子类型 支持数组，参数请参考一言文档 格式x,x  j为网易云，咱这是给女朋友写的，首先排除网易云
    types= a,b,c,d,e,f,g,h,i,k,l
    
    [Day]
    #距离多少天，格式mm-dd，可数组，以,分割。
    count-down=
    #对应距离多少天的的title，顺序对应count-down，可数组，完全描述 （你的描述）x天 如 "距离宝贝的生日还有"
    count-down-title =
    
    #模式 公元/农历 值为 AD/CC （本来想搞个农历，但是处理农历太浪费时间了，实际上农历已经处理完了，但是有个小瑕疵需要我手写农历转换逻辑，不想写，哪天有哪位大佬写的库我腆着脸去伸手）
    #count-down-mode = AD
    
    #迄今为止天数，格式yyyy-mm-dd 可数组 以,分割。
    count =
    #与迄今为止的天数对应的title，顺序对应count，可数组，完全描述 （你的描述）x天 如 "今年已经过了"
    count-title =
    
    #模式 公元/农历 值为 AD/CC （本来想搞个农历，但是处理农历太浪费时间了，实际上农历已经处理完了，但是有个小瑕疵需要我手写农历转换逻辑，不想写，哪天有哪位大佬写的库我去腆着脸伸手）
    #count-mode =
    #是否润年 与count对应 （本来想搞个农历，但是处理农历太浪费时间了，实际上农历已经处理完了，但是有个小瑕疵需要我手写农历转换逻辑，不想写，哪天有哪位大佬写的库我去腆着脸伸手）
    #count-if-lep-year =
    
    [News]
    #是否开始热搜推送（这个我也没写）
    switch = false


消息模板格式：

    {{now.DATA}}

    {{weather1.DATA}}
    
    {{weather2.DATA}}
    
    {{weather3.DATA}}
    
    {{count1.DATA}}
    
    {{count2.DATA}}
    
    {{countDown1.DATA}}
    
    {{countDown2.DATA}}
    
    {{hitokoto.DATA}}

消息模板说明：
    
    {{now.DATA}}    //时间，这个不用动

    {{weather1.DATA}}   //天气，可以看到1，2，3这样的序号，配置文件里可以多配置，相应的在这里写同样数量的配置对应配置文件中的数组顺序就行了。
    
    {{weather2.DATA}}
    
    {{weather3.DATA}}
    
    {{count1.DATA}}     //迄今为止的天数，可以看到有1，2这样的序号，配置文件里可以多配置，相应的在这里写同样数量的配置就对应配置文件中的数组顺序就行了。
    
    {{count2.DATA}}     
    
    {{countDown1.DATA}} //纪念日天数倒计时，可以看到有1，2这样的序号，配置文件里可以多配置，相应的在这里写同样数量的配置就对应配置文件中的数组顺序就行了。
    
    {{countDown2.DATA}}
    
    {{hitokoto.DATA}}   //每日一句，这个不用动


注意事项：
    
    1、配置文件中的每一项都要严格按照我注释中的规定填写。

    2、代码每一块我都分开了，很容易自己改造。

    3、我一个nil都没处理，你打我啊。


开源类库及第三方产品使用：
 - hitokoto(一言) [https://hitokoto.cn/](https://hitokoto.cn/)
 - 日期处理使用了carbon [https://github.com/golang-module/carbon](https://github.com/golang-module/carbon)
 - json处理使用了gjson [https://github.com/tidwall/gjson](https://github.com/tidwall/gjson)
