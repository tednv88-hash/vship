-- Migration 006: Seed default page designs for admin tenant
-- Creates a default homepage (商城首页) for the existing admin tenant
-- New user registrations will get default pages automatically via Go code

-- Insert default homepage for admin tenant (432ba247-0944-40e8-a94f-e58603859020)
-- Only insert if no default homepage exists for this tenant
INSERT INTO page_designs (id, tenant_id, name, type, page_data, status, is_default, created_at, updated_at, extra_fields)
SELECT
    gen_random_uuid(),
    '432ba247-0944-40e8-a94f-e58603859020'::uuid,
    '商城首页',
    'home',
    '{"page":{"type":"page","name":"页面设置","params":{"name":"商城首页","title":"集运商城","share_title":"集运商城"},"style":{"titleTextColor":"white","titleBackgroundColor":"#05ce78"},"id":"page"},"items":[{"name":"图片轮播","type":"banner","style":{"btnColor":"#ffffff","btnShape":"round"},"params":{"interval":"2800"},"data":[{"imgUrl":"/static/img/diy/banner/01.png","linkUrl":""},{"imgUrl":"/static/img/diy/banner/02.png","linkUrl":""}]},{"name":"头条快报","type":"special","params":{"source":"auto","auto":{"category":0,"showNum":6}},"style":{"display":"1","image":"/static/img/diy/special.png"},"defaultData":[{"article_title":"欢迎使用集运商城系统"},{"article_title":"跨境集运，安全快捷"}],"data":[]},{"name":"导航组","type":"navBar","style":{"background":"#ffffff","rowsNum":"4"},"data":[{"imgUrl":"/static/img/diy/navbar/01.png","imgName":"icon-1.png","linkUrl":"","text":"仓库滞留","color":"#666666"},{"imgUrl":"/static/img/diy/navbar/02.png","imgName":"icon-2.png","linkUrl":"","text":"快递盲盒","color":"#666666"},{"imgUrl":"/static/img/diy/navbar/03.png","imgName":"icon-3.png","linkUrl":"","text":"国内特产","color":"#666666"},{"imgUrl":"/static/img/diy/navbar/04.png","imgName":"icon-4.png","linkUrl":"","text":"国外特产","color":"#666666"}]},{"name":"商品组","type":"goods","params":{"source":"auto","auto":{"category":0,"goodsSort":"all","showNum":6}},"style":{"background":"#F6F6F6","display":"list","column":"2","show":{"goodsName":"1","goodsPrice":"1","linePrice":"1","sellingPoint":"0","goodsSales":"0"}},"defaultData":[{"goods_name":"此处显示商品名称","image":"/static/img/diy/goods/01.png","goods_price":"99.00","line_price":"139.00","selling_point":"此款商品美观大方 不容错过","goods_sales":"100"},{"goods_name":"此处显示商品名称","image":"/static/img/diy/goods/01.png","goods_price":"99.00","line_price":"139.00","selling_point":"此款商品美观大方 不容错过","goods_sales":"100"}],"data":[]}]}'::jsonb,
    'active',
    true,
    NOW(),
    NOW(),
    '{}'::jsonb
WHERE NOT EXISTS (
    SELECT 1 FROM page_designs
    WHERE tenant_id = '432ba247-0944-40e8-a94f-e58603859020'
      AND type = 'home'
      AND is_default = true
      AND trashed_at IS NULL
);
