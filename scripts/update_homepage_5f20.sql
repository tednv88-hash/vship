UPDATE page_designs
SET
  page_data = $$
{
  "page": {
    "type": "page",
    "name": "页面设置",
    "params": {
      "name": "商城首页",
      "title": "大發貨運",
      "share_title": "大發貨運"
    },
    "style": {
      "titleTextColor": "white",
      "titleBackgroundColor": "#0f3a57"
    },
    "id": "page"
  },
  "items": [
    {
      "name": "图片轮播",
      "type": "banner",
      "style": {"btnColor": "#d6a93a", "btnShape": "round"},
      "params": {"interval": "2800"},
      "data": [
        {"imgUrl": "/static/img/diy/banner/01.png", "linkUrl": ""},
        {"imgUrl": "/static/img/diy/banner/02.png", "linkUrl": ""}
      ]
    },
    {
      "name": "公告组",
      "type": "notice",
      "params": {"text": "春節放假通知", "icon": "/static/img/diy/notice.png", "moreText": "更多"},
      "style": {"paddingTop": "8", "background": "#ffffff", "textColor": "#333333"}
    },
    {
      "name": "导航组",
      "type": "navBar",
      "style": {"background": "#ffffff", "rowsNum": "5"},
      "data": [
        {"imgUrl": "/static/img/diy/navbar/01.png", "imgName": "icon-1.png", "linkUrl": "", "text": "代購錢包", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/02.png", "imgName": "icon-2.png", "linkUrl": "", "text": "1688代採購", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/03.png", "imgName": "icon-3.png", "linkUrl": "", "text": "集運包裹", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/04.png", "imgName": "icon-4.png", "linkUrl": "", "text": "集運訂單", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/05.png", "imgName": "icon-5.png", "linkUrl": "", "text": "運費充值", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/06.png", "imgName": "icon-6.png", "linkUrl": "", "text": "新手指引", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/07.png", "imgName": "icon-7.png", "linkUrl": "", "text": "倉庫地址", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/08.png", "imgName": "icon-8.png", "linkUrl": "", "text": "運費計算", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/09.png", "imgName": "icon-9.png", "linkUrl": "", "text": "優惠券", "color": "#333333"},
        {"imgUrl": "/static/img/diy/navbar/10.png", "imgName": "icon-10.png", "linkUrl": "", "text": "更多功能", "color": "#333333"}
      ]
    },
    {
      "name": "单图组",
      "type": "imageSingle",
      "style": {"paddingTop": 8, "paddingLeft": 8, "background": "#ffffff"},
      "data": [
        {"imgUrl": "/static/img/diy/window/01.jpg", "imgName": "ad-1.jpg", "linkUrl": ""}
      ]
    },
    {
      "name": "商品组",
      "type": "goods",
      "params": {"source": "auto", "auto": {"category": 0, "goodsSort": "all", "showNum": 8}},
      "style": {
        "background": "#f5f6f7",
        "display": "list",
        "column": "2",
        "show": {"goodsName": "1", "goodsPrice": "1", "linePrice": "0", "sellingPoint": "0", "goodsSales": "0"}
      },
      "defaultData": [
        {"goods_name": "抖音抖加上熱門抖加dou", "image": "/static/img/diy/window/01.jpg", "goods_price": "4.06", "line_price": "0", "selling_point": "", "goods_sales": "0"},
        {"goods_name": "10抖音幣充值速到賬", "image": "/static/img/diy/window/02.jpg", "goods_price": "0.80", "line_price": "0", "selling_point": "", "goods_sales": "0"},
        {"goods_name": "新老顧客不限新老", "image": "/static/img/diy/window/03.jpg", "goods_price": "1.03", "line_price": "0", "selling_point": "", "goods_sales": "0"},
        {"goods_name": "拼豆補充包diy袋裝221色", "image": "/static/img/diy/window/04.jpg", "goods_price": "10.37", "line_price": "0", "selling_point": "", "goods_sales": "0"}
      ],
      "data": []
    }
  ]
}
$$::jsonb,
  updated_at = NOW()
WHERE id = '5f20e844-ae0f-4238-8572-0de2641af208';
