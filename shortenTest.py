import concurrent.futures
import json
import requests
import time
from pathlib import Path

# 读取链接文件
urls_path = Path("test/urls.txt")
with urls_path.open("r", encoding="utf-8") as f:
    urls = [line.strip() for line in f if line.strip()]

# 保存成功的短链接
short_urls = set()
results = []
TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoidXNlciIsImV4cCI6MTc0NDIxOTIzOH0.L3Dma9_iAoC6NoN5JM_LAbnvKlx6nY2xr4JBJmXo4I0"
# 并发发送请求
def send_request(long_url):
    try:
        response = requests.post(
            "http://localhost:8080/api/shorten",
            headers={"Content-Type": "application/json",
                     "Authorization": f"Bearer {TOKEN}"},
            data=json.dumps({"original_url": long_url}),
            timeout=5,
        )
        if response.status_code == 200:
            data = response.json()
            short_url = data.get("shortlink", "null")
            short_urls.add(short_url)
            return f"✅ 成功: {long_url} → {short_url}"
        else:
            return f"❌ 失败: {long_url} → 状态码 {response.status_code}"
    except Exception as e:
        return f"💥 异常: {long_url} → {e}"

start_time = time.time()
# 使用线程池并发请求
with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
    futures = [executor.submit(send_request, url) for url in urls]
    for future in concurrent.futures.as_completed(futures):
        results.append(future.result())

# 结束时间
end_time = time.time()
elapsed = end_time - start_time

# 保存结果日志
log_path = Path("test/results.txt")
log_path.write_text("\n".join(results), encoding="utf-8")



# 输出汇总
print("—— 测试完成 ——")
print(f"总请求数: {len(urls)}")
print(f"成功短链数: {len(short_urls)}")
print(f"唯一短链数: {len(set(short_urls))}")
print(f"总耗时: {elapsed:.2f} 秒")
print(f"平均 QPS: {len(urls)/elapsed:.2f}")
print(f"详细结果已保存到: {log_path}")
