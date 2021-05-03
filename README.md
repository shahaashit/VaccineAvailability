# VaccineAvailability

Periodically checks for availability of vaccine for the pincode specified in config and sends mail if there are any available centers.

### Steps to run app on you machine: 
1. Clone/Download this repo on you machine.
2. Copy config/sample_config.yaml and put it at the same location with name "config.yaml"
3. Update the config.yaml file with the pincode (can put multiple pincodes here with comma separated values eg. "400045,400344,400004") and age for which you want to check availability. Update the email id and password from which you want to send the mail.
4. Update `cronfreq` in config.yaml with values such as "1m", "1h", "30s" to set the frequency at which to query for available slots.
3. Open terminal and run command `./VaccineAvailability`.
4. As its a cron, keep it running (you can minimize the terminal window) and when you want to stop it, press `ctrl+c` from your keyboard in terminal. 

### Note 
1. It will be stopped automatically if you shut down your machine.
2. After switching on the machine, you will have to follow step #3 from above to start the service again.
3. In case any error is printed, please confirm if your internet connection is working fine.