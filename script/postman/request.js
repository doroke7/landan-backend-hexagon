// 1. 直接從 NPM 載入 node-forge 套件，不需要 eval
const forge = pm.require("npm:node-forge");
const CryptoJS = require('crypto-js');

// 2. 定義允許的自訂字元集
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^*()_+-={}|[]";

function generateRandomString(length) {
    let result = "";
    const charsetLength = charset.length;
    for (let i = 0; i < length; i++) {
        const randomValue = CryptoJS.lib.WordArray.random(4).words[0];
        const randomIndex = Math.abs(randomValue) % charsetLength;
        result += charset.charAt(randomIndex);
    }
    return result;
}

// 【新增】將標準 Base64 轉換為 Base64URL 格式的函式
function toBase64Url(base64Str) {
    return base64Str
        .replace(/\+/g, '-') // 把 + 換成 -
        .replace(/\//g, '_') // 把 / 換成 _
        .replace(/=/g, '');  // 移除 =
}

// 3. 隨機產生 16 個字元的 AES Key 與 IV
const randomKeyPlain = generateRandomString(16);
const randomIvPlain = generateRandomString(16);

let oKeys = {
    "key": randomKeyPlain,
    "iv": randomIvPlain,
};
let sKeys = JSON.stringify(oKeys);

const key = CryptoJS.enc.Utf8.parse(randomKeyPlain);
const iv = CryptoJS.enc.Utf8.parse(randomIvPlain);

// 4. 從 Postman 的 Body 中取出原始 JSON 數據
let jsonString = "";
if (pm.request.body && pm.request.body.raw) {
    jsonString = pm.request.body.raw;
} else {
    jsonString = JSON.stringify({ error: "Body was empty" });
}

// 5. 進行 AES-128-CBC 加密
const encrypted = CryptoJS.AES.encrypt(jsonString, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
});

// 這裡將 AES 密文轉為 Base64URL 格式
const sP = toBase64Url(encrypted.toString());

// =================【RSA 加密 sKeys】=================
const rsaPublicKeyPem = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqBbztPHxvEsnb5BcMTjv
693XqRMYte+ORJvrgc0RsdHlC4W4lSqvjnH2JcsTGgtSqmmvsvoPQ9dKGs6+OD3E
zIClDiEr4n7QRFFjYKP3IBqhkR5a5wZdiOCoYCx2dOKjBTkLgIzMO145nITHR0za
Yv7k22eNdIzlLVat1Oq1DlWCHWBEQHUUm/OhiBSHnRb2DXiMa+vBvHHrZBIcDb0+
TRD14zLArY5ijKWkzTLGzr4IDi3TcwDz6xEkLm4grzi/KEYtjAweVTClqm19vYAk
SDe+BtVYNxODv3yQSSIrDEzeCnbimIBCBfwxL65YrbIAUx7YqVbtNry56C4MI95h
rQIDAQAB
-----END PUBLIC KEY-----`;

let sK = "RSA_ENCRYPTION_FAILED";
try {
    const publicKey = forge.pki.publicKeyFromPem(rsaPublicKeyPem);
    const bytesToEncrypt = forge.util.encodeUtf8(sKeys);
    const rsaEncryptedBytes = publicKey.encrypt(bytesToEncrypt, 'RSAES-PKCS1-V1_5');
    
    // 這裡將 RSA 密文轉為 Base64URL 格式
    // base64 包含 url 特殊字元，不能用一般的 base64編碼
    // base64 包含 url 特殊字元，不能用一般的 base64編碼
    // base64 包含 url 特殊字元，不能用一般的 base64編碼
    // base64 包含 url 特殊字元，不能用一般的 base64編碼
    // base64 包含 url 特殊字元，不能用一般的 base64編碼
    // base64 包含 url 特殊字元，不能用一般的 base64編碼

    sK = forge.util.encode64(rsaEncryptedBytes);
    sK = toBase64Url(sK);

    pm.variables.set("K", sK);
} catch (e) {
    console.error("RSA 加密失敗：", e.message);
}

pm.variables.set("p", sP);

// =================【重寫 Body 只發送 {"p": sP}】=================
const oBody = {
    "p": sP
};

pm.request.body.update({
    mode: 'raw',
    raw: JSON.stringify(oBody)
});

// 強制將 Content-Type 設定為 application/json
pm.request.headers.upsert({
    key: "Content-Type",
    value: "application/json"
});

// =================【把 sK 塞進 Header】=================
pm.request.headers.upsert({
    key: "K",
    value: sK
});
// =====================================================

console.log("--- 【混合加密 (Base64URL 格式) 成功】 ---");
console.log("已覆蓋發送的 Body (Base64URL)：", JSON.stringify(oBody));
console.log("已新增 Header [K] (Base64URL)：", sK);