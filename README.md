# VaccineAvailability

Periodically checks for availability of vaccine for the pincode specified in config.

### Steps to run app on you machine: 
1. Clone/Download this repo on you machine.
2. Update the config.yaml file with the pincode (can put multiple pincodes here with comma separated values eg. "400045,400344,400004") and age for which you want to check availability.
3. Open terminal and run command `./VaccineAvailability`.
4. As its a cron, keep it running (you can minimize the terminal window) and when you want to stop it, press `ctrl+c` from your keyboard in terminal. 

### Note 
1. Availability will be checked every minute.
2. It will be stopped automatically if you shut down your machine.
3. After switching on the machine, you will have to follow step #3 from above to start the service again.
4. If any center has availability for a given configuration, it will be printed in terminal. (Notification over email is WIP).
5. In case any error is printed, please confirm if your internet connection is working fine.