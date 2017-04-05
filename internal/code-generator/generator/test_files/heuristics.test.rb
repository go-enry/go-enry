    # Common heuristics
    ObjectiveCRegex = /^\s*(@(interface|class|protocol|property|end|synchronised|selector|implementation)\b|#import\s+.+\.h[">])/

    disambiguate ".asc" do |data|
      if /^(----[- ]BEGIN|ssh-(rsa|dss)) /.match(data)
        Language["Public Key"]
      elsif /^[=-]+(\s|\n)|{{[A-Za-z]/.match(data)
        Language["AsciiDoc"]
      elsif /^(\/\/.+|((import|export)\s+)?(function|int|float|char)\s+((room|repeatedly|on|game)_)?([A-Za-z]+[A-Za-z_0-9]+)\s*[;\(])/.match(data)
        Language["AGS Script"]
      end
    end

    disambiguate ".ms" do |data|
      if /^[.'][a-z][a-z](\s|$)/i.match(data)
        Language["Groff"]
      elsif /(?<!\S)\.(include|globa?l)\s/.match(data) || /(?<!\/\*)(\A|\n)\s*\.[A-Za-z]/.match(data.gsub(/"([^\\"]|\\.)*"|'([^\\']|\\.)*'|\\\s*(?:--.*)?\n/, ""))
        Language["GAS"]
      else
        Language["MAXScript"]
      end
    end

    disambiguate ".mod" do |data|
      if data.include?('<!ENTITY ')
        Language["XML"]
      elsif /^\s*MODULE [\w\.]+;/i.match(data) || /^\s*END [\w\.]+;/i.match(data)
        Language["Modula-2"]
      else
        [Language["Linux Kernel Module"], Language["AMPL"]]
      end
    end

    disambiguate ".pro" do |data|
      if /^[^#]+:-/.match(data)
        Language["Prolog"]
      elsif data.include?("last_client=")
        Language["INI"]
      elsif data.include?("HEADERS") && data.include?("SOURCES")
        Language["QMake"]
      elsif /^\s*function[ \w,]+$/.match(data)
        Language["IDL"]
      end
    end
