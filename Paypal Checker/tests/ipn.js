var tap = require('tap')
var qs = require('querystring')
var ipn = require('../ipn.js')
var fixtures = require('./fixtures.js')

// Test helper functions
tap.equal(ipn.stringify({ one: '1', two: '2', three: '3' }), '&one=1&two=2&three=3')
tap.equal(ipn.stringify({ one: '%1', two: '2', three: '3' }), '&one=%1&two=2&three=3')
tap.same(ipn.parseQuery('?mc_gross=19.95'), { 'mc_gross': 19.95 })
tap.notSame(ipn.parseQuery('?mc_gross=19.00'), { 'mc_gross': 19.95 })

tap.test('fake ipn request', function (childTest) {
  var params = qs.parse(fixtures.invalidIpn())
  ipn.verify(params, { 'allow_sandbox': true }, function callback (err) {
    childTest.match(err, /Error/)
    childTest.end()
  })
})

tap.test('valid ipn request', function (childTest) {
  var params = qs.parse(fixtures.validIpn())
  ipn.verify(params, { 'allow_sandbox': true }, function callback (err) {
    childTest.match(err, null)
    childTest.end()
  })
})

tap.test('valid test ipn request but sandbox not allowed', function (childTest) {
  var params = qs.parse(fixtures.validIpn())
  ipn.verify(params, { 'allow_sandbox': false }, function callback (err) {
    childTest.match(err, /Received request with test_ipn parameter while sandbox is disabled/)
    childTest.end()
  })
})

tap.test('valid ipn request with diacritics in item_name', function (childTest) {
  ipn.verify(fixtures.validIpnObject01(), { 'allow_sandbox': true }, function callback (err) {
    childTest.match(err, null)
    childTest.end()
  })
})

tap.test('test without params', function (childTest) {
  ipn.verify(undefined, { 'allow_sandbox': true }, function callback (err) {
    childTest.match(err, /No params were passed to ipn.verify/)
    childTest.end()
  })
})

tap.test('test without options', function (childTest) {
  ipn.verify(fixtures.validIpnObject01(), function callback (err) {
    childTest.match(err, /Received request with test_ipn parameter while sandbox is disabled/)
    childTest.end()
  })
})