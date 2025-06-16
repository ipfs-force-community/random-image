from diffusers import StableDiffusionPipeline
import torch
import os
import time
import random
import string
import argparse

# 解析命令行参数
parser = argparse.ArgumentParser(description="Generate images with Stable Diffusion")
parser.add_argument("--output_dir", type=str, default="generated_images", help="Output directory for images")
parser.add_argument("--num_images", type=int, default=5, help="Number of images to generate")
parser.add_argument("--prompt", type=str, default="一幅赛博朋克风格的城市夜景，霓虹灯，未来感", help="Prompt for image generation")
args = parser.parse_args()

# 创建保存图片的文件夹
output_dir = args.output_dir
os.makedirs(output_dir, exist_ok=True)

# 加载模型
pipe = StableDiffusionPipeline.from_pretrained(
    "runwayml/stable-diffusion-v1-5",
    torch_dtype=torch.float16,
).to("cuda")

# 生成 14 位小写字母或数字的随机文件名
def generate_random_filename():
    chars = string.ascii_lowercase + string.digits  # a-z, 0-9
    return ''.join(random.choice(chars) for _ in range(14)) + '_' + ''.join(random.choice(chars) for _ in range(14)) + ".png"

# 生成图片并统计时间
prompt = args.prompt
num_images = args.num_images  # 生成图片数量

def generate_random_image():
    print(f"开始生成 {num_images} 张图片，保存到 {output_dir} 目录...")
    print(f"使用的提示词: {prompt}")
    for i in range(num_images):
        # 记录开始时间
        start_time = time.time()

        # 生成图片
        image = pipe(prompt, num_inference_steps=50, guidance_scale=7.5).images[0]

        # 生成随机文件名
        filename = generate_random_filename()
        filepath = os.path.join(output_dir, filename)

        # 保存调整后的图片
        image.save(filepath, format="PNG", optimize=True)
        
        # 记录结束时间并计算耗时
        end_time = time.time()
        elapsed_time = end_time - start_time

        # 打印信息
        print(f"图片 {i+1}: {filename}，生成时间: {elapsed_time:.2f} 秒")

if __name__ == "__main__":
    print(generate_random_image())
