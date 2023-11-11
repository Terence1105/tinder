### 系統設計文檔

#### 概述
這個系統是一個基於 HTTP API 的配對系統，用於模擬類似於 Tinder 的功能。主要目標是讓用戶加入系統、尋找配對，以及從系統中移除用戶。

#### API

1. **加入用戶並配對 (`/v1/add-single-person-and-match`)**
   - **方法**: POST
   - **功能**: 加入一個新用戶到系統中，並嘗試為這個用戶找到配對。
   - **請求參數**: `name` (string), `height` (float64), `gender` (int), `dateCounts` (int)

2. **移除用戶 (`/v1/remove-single-person`)**
   - **方法**: POST
   - **功能**: 從系統中移除一個用戶，使其不再參與配對。
   - **請求參數**: `name` (string), `gender` (int)

3. **查詢配對 (`/v1/query-single-people`)**
   - **方法**: GET
   - **功能**: 查詢系統中最多 `N` 組可能的配對。
   - **請求參數**: `counts` (int)

#### 數據存儲
使用 Redis 作為後端存儲，主要存儲用戶資料和配對資訊。

#### 時間複雜度分析

1. **加入用戶並配對**
   - 時間複雜度：O(1)
   - 說明：添加用戶到 Redis 是一個常數時間操作，不依賴於用戶數量。

2. **移除用戶**
   - 時間複雜度：O(1)
   - 說明：從 Redis 中移除用戶同樣是一個常數時間操作。

3. **查詢配對**
   - 時間複雜度：O(N*M)
   - 說明：對於每個男性用戶（假設有 N 個），系統可能需要查詢多達 M 個女性用戶（這取決於每個男性用戶的約會次數和女性用戶的可用性）。因此，時間複雜度與男性用戶數量和其約會次數成正比。

#### 可在本地端 `make run` 運行，API文件請參考 `make run-swag`
