# from docx import Document
# import re


# # doc = Document(r'C:\Users\23725\Downloads\生理答案（全部）.docx')
# # pattern = re.compile(r'(正确答案： )([A-Za-z]+)')
# # for para in doc.paragraphs:
# #     para.text = pattern.sub(r'\1', para.text)
# # doc.save(r'C:\Users\23725\Downloads\生理答案（全部）1.docx')

# # doc = Document(r'C:\Users\23725\Downloads\生理答案（全部）.docx')
# # for para in doc.paragraphs:
# #     match = re.search(r'(正确答案： )([A-Za-z]+)', para.text)
# #     if match:
# #         para.text = para.text.replace(match.group(0), match.group(1))
# # doc.save(r'C:\Users\23725\Downloads\生理答案（全部）1.docx')

# doc = Document(r'C:\Users\23725\Downloads\生理答案（全部）.docx')
# for para in doc.paragraphs:
#     for run in para.runs:
#         run.text = re.sub(r'[：:]\s*[A-Za-z]+', '：', run.text)
# doc.save(r'C:\Users\23725\Downloads\生理答案（全部）1.docx')

# import hashlib
# from ecdsa import SECP256k1, VerifyingKey, SigningKey

# # 生成密钥对（私钥和公钥）
# def generate_key_pair():
#     # 创建一个SigningKey对象（私钥）
#     signing_key = SigningKey.generate(curve=SECP256k1)
#     private_key = signing_key.to_string().hex()

#     # 通过SigningKey对象获取公钥（公钥是私钥的椭圆曲线签名的公钥）
#     public_key = signing_key.get_verifying_key().to_string().hex()

#     return private_key, public_key

# # 计算公钥的哈希（RIPEMD160(SHA256(pubkey))）
# def public_key_hash(public_key):
#     # 将公钥转换为字节
#     public_key_bytes = bytes.fromhex(public_key)

#     # 使用SHA256对公钥进行哈希
#     sha256_hash = hashlib.sha256(public_key_bytes).digest()

#     # 使用RIPEMD160对SHA256哈希结果进行哈希
#     ripemd160_hash = hashlib.new('ripemd160', sha256_hash).digest()

#     return ripemd160_hash.hex()

# # 生成签名
# def sign_message(private_key, message):
#     # 使用私钥生成签名
#     signing_key = SigningKey.from_string(bytes.fromhex(private_key), curve=SECP256k1)
#     signature = signing_key.sign(message.encode())

#     return signature.hex()

# # 验证签名
# def verify_signature(public_key, message, signature):
#     # 将签名转换为字节
#     signature_bytes = bytes.fromhex(signature)

#     # 将公钥转换为VerifyingKey对象
#     public_key_bytes = bytes.fromhex(public_key)
#     verifying_key = VerifyingKey.from_string(public_key_bytes, curve=SECP256k1)

#     # 验证签名
#     try:
#         verifying_key.verify(signature_bytes, message.encode())
#         return True  # 签名有效
#     except:
#         return False  # 签名无效

# # 示例流程
# private_key, public_key = generate_key_pair()
# print(f"Private Key: {private_key}")
# print(f"Public Key: {public_key}")

# # 获取公钥哈希
# public_key_hash_value = public_key_hash(public_key)
# print(f"Public Key Hash: {public_key_hash_value}")

# # 签名
# message = "Hello, Bitcoin!"
# signature = sign_message(private_key, message)
# print(f"Signature: {signature}")

# # 验证签名
# is_valid = verify_signature(public_key, message, signature)
# print(f"Signature valid: {is_valid}")


import hashlib

pwd = '123'
salt = 'e97e4454e5f3df8b4c571eb0fd420b0d'
hash = hashlib.sha256((pwd + salt).encode()).hexdigest()
print(hash)
