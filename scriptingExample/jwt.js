/**
 * 生成 JWT Token
 * @author ClipSave
 * @param {string} memberCode - 会员码，从剪贴板内容获取
 * @param {string} secretKey - 密钥
 * @returns {string} - 生成的 JWT Token
 * @returns {object} - 错误信息
 */

const secretKey = "your-secret-key";

// 主逻辑：从剪贴板内容获取 memberCode
const memberCode = item.Content ? item.Content.trim() : "";

if (!memberCode) {
  return {
    error: "剪贴板内容为空，无法生成 token"
  };
}

function base64UrlEncode(str) {
  const encoder = new TextEncoder();
  const bytes = encoder.encode(str);
  const binary = String.fromCharCode(...bytes);
  return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, "");
}

async function generateJWT(memberCode, secretKey) {
  const header = { alg: "HS256" };
  const now = Math.floor(Date.now() / 1000);

  const claims = {
    exp: now + 6 * 60 * 60, // 6小时后过期
    iat: now,                // 签发时间
    memberCode: memberCode,  // memberCode 在最后
  };

  const encodedHeader = base64UrlEncode(JSON.stringify(header));
  const encodedPayload = base64UrlEncode(JSON.stringify(claims));
  const signatureInput = `${encodedHeader}.${encodedPayload}`;

  let secretBytes;
  try {
    const decodedSecret = atob(secretKey);
    secretBytes = new Uint8Array(decodedSecret.length);
    for (let i = 0; i < decodedSecret.length; i++) {
      secretBytes[i] = decodedSecret.charCodeAt(i);
    }
  } catch (e) {
    const encoder = new TextEncoder();
    secretBytes = encoder.encode(secretKey);
  }

  // 使用解码后的密钥生成签名
  const key = await crypto.subtle.importKey(
    "raw",
    secretBytes,
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"]
  );

  const encoder = new TextEncoder();
  const signature = await crypto.subtle.sign(
    "HMAC",
    key,
    encoder.encode(signatureInput)
  );
  
  const signatureArray = Array.from(new Uint8Array(signature));
  const signatureBase64 = btoa(String.fromCharCode(...signatureArray))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");

  return `${encodedHeader}.${encodedPayload}.${signatureBase64}`;
}

// 生成 token
const token = await generateJWT(memberCode, secretKey);

return token;
