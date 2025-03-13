-- 设置请求方法为 POST
wrk.method = "POST"

-- 设置请求头
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNzc1NTk5ODExNDQ4MDEyOSwiZXhwIjoxNzQxODkxMTAwLCJpc3MiOiJDQ01PIn0.wWVNVd9uKr6iBAq9VHy1uDRk7JZVDSeq5jfdGsQQefM"

-- 设置请求体数据（JSON格式）
wrk.body = '{}'

-- 返回请求体
function request()
    return wrk.format(nil, nil, nil, wrk.body)
end