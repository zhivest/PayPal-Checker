module.exports = {
  validIpnObject01: function () {
    return {
      'mc_gross': '0.50',
      'protection_eligibility': 'Ineligible',
      'payer_id': 'J9VXV8D7RGTDU',
      'tax': '0.00',
      'payment_date': '02:37:16 Feb 10, 2016 PST',
      'payment_status': 'Completed',
      'charset': 'windows-1252',
      'first_name': 'Test',
      'mc_fee': '0.32',
      'notify_version': '3.8',
      'custom': '$2a$10$nHMVxF1jBOoeLJlpS.1kiuJlyP2vIqfae4glTSmWUK6M7OOT33cku',
      'payer_status': 'verified',
      'business': 'paypal@oronjo.nl',
      'quantity': '1',
      'verify_sign': 'AGs52InFHk3cJ5j80u2Udq.A26J4AF-BoXrl2vCiBrcjmBKzSdSKDohx',
      'payer_email': 'ronald.luitwieler-buyer@gmail.com',
      'txn_id': '5RJ61844Y3467971C',
      'payment_type': 'instant',
      'last_name': 'Buyer',
      'receiver_email': 'paypal@oronjo.nl',
      'payment_fee': '0.32',
      'receiver_id': 'TCSWDLLDD6Q7W',
      'txn_type': 'web_accept',
      'item_name': '+%E0+pr%F6duct',
      'mc_currency': 'USD',
      'item_number': 'GpaFdqvDW6LN8y9Xg',
      'residence_country': 'NL',
      'test_ipn': '1',
      'handling_amount': '0.00',
      'transaction_subject': '$2a$10$nHMVxF1jBOoeLJlpS.1kiuJlyP2vIqfae4glTSmWUK6M7OOT33cku',
      'payment_gross': '0.50',
      'shipping': '0.00',
      'ipn_track_id': '301ff34d9a01'
    }
  },

  invalidIpn: function () {
    return 'mc_gross=19.95&protection_eligibility=Eligible&address_status=confirmed&payer_id=LPLWNMTBWMFAY&tax=0.00&address_street=1+Main+St&payment_date=20%3A12%3A59+Jan+13%2C+2009+PST&payment_status=Completed&charset=windows-1252&address_zip=95131&first_name=Test&mc_fee=0.88&address_country_code=US&address_name=Test+User&notify_version=2.6&custom=&payer_status=verified&address_country=United+States&address_city=San+Jose&quantity=1&verify_sign=AtkOfCXbDm2hu0ZELryHFjY-Vb7PAUvS6nMXgysbElEn9v-1XcmSoGtf&payer_email=gpmac_1231902590_per%40paypal.com&txn_id=61E67681CH3238416&payment_type=instant&last_name=User&address_state=CA&receiver_email=gpmac_1231902686_biz%40paypal.com&payment_fee=0.88&receiver_id=S8XGHLYDW9T3S&txn_type=express_checkout&item_name=&mc_currency=USD&item_number=&residence_country=US&test_ipn=1&handling_amount=0.00&transaction_subject=&payment_gross=19.95&shipping=0.00'
  },

  validIpn: function () {
    return 'cmd=_notify-validate&mc_gross=32.00&protection_eligibility=Ineligible&payer_id=UDDTGF8DURN4Q&tax=0.00&payment_date=04%3A22%3A08+Oct+23%2C+2013+PDT&payment_status=Completed&charset=windows-1252&first_name=ronald&mc_fee=1.55&notify_version=3.7&custom=191c812474ea53c4646a4c3495bc3b1979a38cbb%3A15fb0c237cca8dcc4338137e94edb46d&payer_status=verified&business=oronjo1.test.shopper%40orangewise.com&quantity=1&verify_sign=AFcWxV21C7fd0v3bYYYRCpSSRl31A.Nb4StWHTSKyQTE6RS7nJNEPXhx&payer_email=ronald.paypal.tester%40orangewise.com&txn_id=1HM582610N8901528&payment_type=instant&payer_business_name=ronald+paypal+tester+at+orangewise.com&last_name=luitwieler+paypal+tester&receiver_email=oronjo1.test.shopper%40orangewise.com&payment_fee=1.55&receiver_id=UJETUZKE6XD5Q&txn_type=web_accept&item_name=cart+with+5+items&mc_currency=USD&item_number=70772406&residence_country=US&test_ipn=1&handling_amount=0.00&transaction_subject=191c812474ea53c4646a4c3495bc3b1979a38cbb%3A15fb0c237cca8dcc4338137e94edb46d&payment_gross=32.00&shipping=0.00&ipn_track_id=c02dbc986674d'
  }

}