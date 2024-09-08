# brutus


## set up local env

install virtual environment 
```bash
python3 -m venv .venv
```
activate venv
```bash
source .venv/bin/activate
``` 
```fish
source .venv/bin/activate.fish
```
```batch
venv\Scripts\activate.bat
```
```powershell
venv\Scripts\activate.ps1
```
install requirements, in this example for az stuff
```bash
pip install -r ./src/server/az/requirements.txt
``` 