# proto specification

## request

- byte: placeholder
- short: service name
- byte: action name
- int: sequence in the websocket lifecycle, begin from 0 and plus1 for every one
- other bytes: Protobuf

## response

- byte: result metadata data
    - 低1位表示结果成功或者失败: 1 - fail; 0 - success
- short: service name
- byte: action name
- int: sequence
- other bytes: Protobuf

## NOTICE

- Little Endian as bytes sequence
