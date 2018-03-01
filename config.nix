{
  Logger = {
    Level = "info";
    Formatter = "json";
  };
  Inputs = [
    {
      Consumer = {
        Format = "json";
        Queue = {
          Type = "nsq";
          Nsq = {
            Addr = "127.0.0.1:4150";
            Channel = "gluttony";
            ConsumerBufferSize = 128;
            Topic = "tickers";
            LogLevel = "info";
          };
        };
      };

      Preprocessor = {
        Type = "lua";
        Lua = {
          Code = ''
            function preprocess(v)
              v['currencyPair'] = v['currencyPair'][1]['symbol'] ..
                '-' .. v['currencyPair'][2]['symbol']
              return v
            end
          '';
          FunctionName = "preprocess";
          Workers = 128;
        };
      };

      Database = {
        Type = "influxdb";
        Influxdb = {
          Client = {
            Addr = "http://127.0.0.1:8086";
          };
          Writer = {
            Batch = {
              FlushInterval = "5s";
              Points = {
                Precision = "ns";
                Database = "gluttony";
                RetentionPolicy = "autogen";
                WriteConsistency = "any";
              };
              Size = 2048;
            };
            Point = {
              Name = "tickers";
              Fields = [
                "high"
                "low"
                "vol"
                "last"
                "buy"
                "sell"
              ];
              Tags = [
                "market"
                "currencyPair"
                "tags"
              ];
              Timestamp = "timestamp";
              TimestampPrecision = "nanosecond";
            };
          };
        };
      };
    }

  ];
}
