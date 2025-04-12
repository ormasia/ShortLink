在 GORM 中使用 `map[string]any` 作为 `Updates` 方法的参数，主要是为了 **更灵活地控制更新的字段**，尤其是需要处理以下场景时：

---

### 1. **更新零值字段**
   - **问题**：如果使用结构体传递参数，GORM 默认会忽略零值字段（如 `""`, `0`, `false`）。
   - **场景**：若需要将字段显式设置为零值（例如将 `block_reason` 清空为 `""`），必须用 `map`。
   - **对比**：
     ```go
     // 使用结构体（零值字段被忽略）
     updates := URLMapping{Status: status, BlockReason: ""}
     db.Updates(updates) // BlockReason 不会被更新

     // 使用 map（零值字段被更新）
     db.Updates(map[string]any{"status": status, "block_reason": ""}) // BlockReason 被设置为空字符串
     ```

---

### 2. **动态字段选择**
   - **场景**：需要根据条件动态决定更新哪些字段。
   - **示例**：
     ```go
     updates := map[string]any{"status": status}
     if blockReason != "" {
         updates["block_reason"] = blockReason
     }
     db.Updates(updates) // 动态添加字段
     ```

---

### 3. **避免字段名与数据库列名的隐式映射**
   - **问题**：结构体的字段名默认通过驼峰命名法映射到数据库的蛇形列名（如 `BlockReason` → `block_reason`），但可能需要手动通过 `gorm` tag 指定列名。
   - **优势**：使用 `map` 可以直接使用数据库列名，避免结构体与表结构的耦合。
     ```go
     // 结构体需要定义 tag
     type URLMapping struct {
         BlockReason string `gorm:"column:block_reason"`
     }

     // 使用 map 直接指定列名
     db.Updates(map[string]any{"block_reason": "..."}) // 无需依赖结构体定义
     ```

---

### 4. **支持非结构体字段**
   - **场景**：如果需要更新非模型字段（例如关联表字段或计算字段），`map` 更灵活。
   - **示例**：
     ```go
     db.Updates(map[string]any{
         "counter": gorm.Expr("counter + ?", 1), // 直接使用 SQL 表达式
     })
     ```

---

### 总结：为什么用 `map`？
| 场景                  | 结构体              | map               |
|-----------------------|---------------------|-------------------|
| 更新零值字段           | ❌ 忽略零值          | ✅ 支持零值        |
| 动态字段选择           | ❌ 需预先定义        | ✅ 灵活增减字段     |
| 直接指定数据库列名      | ❌ 依赖结构体 tag    | ✅ 无需耦合        |
| 使用 SQL 表达式或函数  | ❌ 不支持           | ✅ 支持            |

在需要 **精确控制更新字段**（尤其是涉及零值或动态逻辑）时，`map` 是更合适的选择。如果只是简单更新非零字段，结构体更简洁。