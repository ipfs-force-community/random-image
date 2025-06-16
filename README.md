# random-image

### build

```
git clone https://github.com/ipfs-force-community/random-image

make
```

### Usage

#### 安装 Python 依赖

```
python3 -m venv env
source ./env/bin/activate

pip install -r requirements.txt
```

> pip install diffusers transformers torch accelerate

#### 生成图片

用于生成 512 * 512 的图片，图片大小大约为 500K

* --output_dir 新生成图片的目录
* --num_images 生成多少张图片
* --prompt 图片提示词，比如：cat，就会生成和猫相关的图片

```
source ./env/bin/activate

python src/main.py --output_dir data --num_images 10 --prompt cat
```

#### 调整图片大小

主要用于把小于 1M 的图片调整成 10M - 50M 的图片，图片会变模糊

* --output 新生成图片的目录
* --source-dir 小图片的目录

```
./random-image g --output xxx --source-dir xxxx
```
