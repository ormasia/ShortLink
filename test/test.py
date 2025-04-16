import concurrent.futures
import json
import requests
import time
from pathlib import Path
from threading import Lock

# 文件路径
urls_path = Path("test/urls.txt")
log_path = Path("test/results.txt")

# 读取链接
with urls_path.open("r", encoding="utf-8") as f:
    urls = [line.strip() for line in f if line.strip()]

# 并发控制 + 连接复用
session = requests.Session()  # ✅ 使用连接池
adapter = requests.adapters.HTTPAdapter(pool_connections=100, pool_maxsize=100)
session.mount('http://', adapter)

# 日志收集
short_urls = set()
results = []
lock = Lock()

TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoidXNlciIsImV4cCI6MTc0NDM0MDA1NH0.Yu8RM1RPKe9osqYke3I3S6BiiZMkmDSQTBLzfK37LbQ"

def send_request(long_url):
    try:
        response = session.post(
            "http://localhost:8080/api/shorten",
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {TOKEN}",
                "Connection": "keep-alive"  # 可选，加固连接复用
            },
            data=json.dumps({"original_url": long_url}),
            timeout=3,
        )
        if response.status_code == 200:
            data = response.json()
            short_url = data.get("shortlink", "null")
            with lock:
                short_urls.add(short_url)
            return f"✅ 成功: {long_url} → {short_url}"
        else:
            return f"❌ 失败: {long_url} → 状态码 {response.status_code}"
    except Exception as e:
        return f"💥 异常: {long_url} → {e}"

# 执行测试
start_time = time.time()
with concurrent.futures.ThreadPoolExecutor(max_workers=200) as executor:  # ✅ 建议从 10~20 开始测试
    futures = [executor.submit(send_request, url) for url in urls]
    for future in concurrent.futures.as_completed(futures):
        results.append(future.result())
end_time = time.time()

# 输出 & 保存
elapsed = end_time - start_time
log_path.write_text("\n".join(results), encoding="utf-8")

# 显示关闭会话
session.close()

print("—— 测试完成 ——")
print(f"总请求数: {len(urls)}")
print(f"成功短链数: {len(short_urls)}")
print(f"唯一短链数: {len(set(short_urls))}")
print(f"总耗时: {elapsed:.2f} 秒")
print(f"平均 QPS: {len(urls)/elapsed:.2f}")
print(f"详细结果已保存到: {log_path}")
