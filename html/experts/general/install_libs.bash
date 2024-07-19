# composer does not support git subdirectories on monorepos 
# (despite 9 years people asking for it)
# But don't worry, there is a workaround - if you use their paid product :D
# 
# So yeah, we do just use bash scripts for installing some stuff here.
# Other stuff installs well with composer so use that in that case.

bash <(curl -s "https://raw.githubusercontent.com/mlaak/robot32lib/main/vendor/Robot32lib/Middleware/install.bash")
bash <(curl -s "https://raw.githubusercontent.com/mlaak/robot32lib/main/vendor/Robot32lib/GPTlib/install.bash")
bash <(curl -s "https://raw.githubusercontent.com/mlaak/robot32lib/main/vendor/Robot32lib/ULogger/install.bash")
bash <(curl -s "https://raw.githubusercontent.com/mlaak/robot32lib/main/vendor/Robot32lib/LLMServerList/install.bash")


