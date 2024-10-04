#!/usr/bin/env bash

# generate mock files
## app/database
mockgen -source=../app/database/jrp/repository/repository.go -destination=../mock/app/database/jrp/repository/repository.go -package=mockrepository
mockgen -source=../app/database/wnjpn/repository/repository.go -destination=../mock/app/database/wnjpn/repository/repository.go -package=mockrepository
## app/library
mockgen -source=../app/library/dbfiledirpathprovider/dbfiledirpathprovider.go -destination=../mock/app/library/dbfiledirpathprovider/dbfiledirpathprovider.go -package=mockdbfiledirpathprovider
mockgen -source=../app/library/downloader/downloader.go -destination=../mock/app/library/downloader/downloader.go -package=mockdownloader
mockgen -source=../app/library/generator/generator.go -destination=../mock/app/library/generator/generator.go -package=mockgenerator
mockgen -source=../app/library/jrpwriter/jrpwriter.go -destination=../mock/app/library/jrpwriter/jrpwriter.go -package=mockjrpwriter
mockgen -source=../app/library/utility/utility.go -destination=../mock/app/library/utility/utility.go -package=mockutility
mockgen -source=../app/library/versionprovider/versionprovider.go -destination=../mock/app/library/versionprovider/versionprovider.go -package=mockversionprovider
## app/proxy
mockgen -source=../app/proxy/buffer/bufferproxy.go -destination=../mock/app/proxy/buffer/bufferproxy.go -package=mockbufferproxy
mockgen -source=../app/proxy/cobra/cobraproxy.go -destination=../mock/app/proxy/cobra/cobraproxy.go -package=mockcobraproxy
mockgen -source=../app/proxy/cobra/commandinstance.go -destination=../mock/app/proxy/cobra/commandinstance.go -package=mockcobraproxy
mockgen -source=../app/proxy/cobra/positionalargsinstance.go -destination=../mock/app/proxy/cobra/positionalargsinstance.go -package=mockcobraproxy
mockgen -source=../app/proxy/color/colorproxy.go -destination=../mock/app/proxy/color/colorproxy.go -package=mockcolorproxy
mockgen -source=../app/proxy/debug/debugproxy.go -destination=../mock/app/proxy/debug/debugproxy.go -package=mockdebugproxy
mockgen -source=../app/proxy/debug/buildinfoinstance.go -destination=../mock/app/proxy/debug/buildinfoinstance.go -package=mockdebugproxy
mockgen -source=../app/proxy/filepath/filepathproxy.go -destination=../mock/app/proxy/filepath/filepathproxy.go -package=mockfilepathproxy
mockgen -source=../app/proxy/fmt/fmtproxy.go -destination=../mock/app/proxy/fmt/fmtproxy.go -package=mockfmtproxy
mockgen -source=../app/proxy/fs/fsproxy.go -destination=../mock/app/proxy/fs/fsproxy.go -package=mockfsproxy
mockgen -source=../app/proxy/fs/fileinfoinstance.go -destination=../mock/app/proxy/fs/fileinfoinstance.go -package=mockfsproxy
mockgen -source=../app/proxy/gzip/gzipproxy.go -destination=../mock/app/proxy/gzip/gzipproxy.go -package=mockgzipproxy
mockgen -source=../app/proxy/http/httpproxy.go -destination=../mock/app/proxy/http/httpproxy.go -package=mockhttpproxy
mockgen -source=../app/proxy/http/responseinstance.go -destination=../mock/app/proxy/http/responseinstance.go -package=mockhttpproxy
mockgen -source=../app/proxy/io/ioproxy.go -destination=../mock/app/proxy/io/ioproxy.go -package=mockioproxy
mockgen -source=../app/proxy/os/osproxy.go -destination=../mock/app/proxy/os/osproxy.go -package=mockosproxy
mockgen -source=../app/proxy/os/fileinstance.go -destination=../mock/app/proxy/os/fileinstance.go -package=mockosproxy
mockgen -source=../app/proxy/pflag/pflagproxy.go -destination=../mock/app/proxy/pflag/pflagproxy.go -package=mockpflagproxy
mockgen -source=../app/proxy/pflag/flagsetinstance.go -destination=../mock/app/proxy/pflag/flagsetinstance.go -package=mockpflagproxy
mockgen -source=../app/proxy/promptui/promptuiproxy.go -destination=../mock/app/proxy/promptui/promptuiproxy.go -package=mockpromptuiproxy
mockgen -source=../app/proxy/promptui/promptinstance.go -destination=../mock/app/proxy/promptui/promptinstance.go -package=mockpromptuiproxy
mockgen -source=../app/proxy/rand/randproxy.go -destination=../mock/app/proxy/rand/randproxy.go -package=mockrandproxy
mockgen -source=../app/proxy/sort/sortproxy.go -destination=../mock/app/proxy/sort/sortproxy.go -package=mocksortproxy
mockgen -source=../app/proxy/spinner/spinnerproxy.go -destination=../mock/app/proxy/spinner/spinnerproxy.go -package=mockspinnerproxy
mockgen -source=../app/proxy/spinner/spinnerinstance.go -destination=../mock/app/proxy/spinner/spinnerinstance.go -package=mockspinnerproxy
mockgen -source=../app/proxy/sql/sqlproxy.go -destination=../mock/app/proxy/sql/sqlproxy.go -package=mocksqlproxy
mockgen -source=../app/proxy/sql/dbinstance.go -destination=../mock/app/proxy/sql/dbinstance.go -package=mocksqlproxy
mockgen -source=../app/proxy/sql/nullstringinstance.go -destination=../mock/app/proxy/sql/nullstringinstance.go -package=mocksqlproxy
mockgen -source=../app/proxy/sql/resultinstance.go -destination=../mock/app/proxy/sql/resultinstance.go -package=mocksqlproxy
mockgen -source=../app/proxy/sql/rowinstance.go -destination=../mock/app/proxy/sql/rowinstance.go -package=mocksqlproxy
mockgen -source=../app/proxy/sql/rowsinstance.go -destination=../mock/app/proxy/sql/rowsinstance.go -package=mocksqlproxy
mockgen -source=../app/proxy/sql/stmtinstance.go -destination=../mock/app/proxy/sql/stmtinstance.go -package=mocksqlproxy
mockgen -source=../app/proxy/sql/txinstance.go -destination=../mock/app/proxy/sql/txinstance.go -package=mocksqlproxy
mockgen -source=../app/proxy/strconv/strconvproxy.go -destination=../mock/app/proxy/strconv/strconvproxy.go -package=mockstrconvproxy
mockgen -source=../app/proxy/strings/stringsproxy.go -destination=../mock/app/proxy/strings/stringsproxy.go -package=mockstringsproxy
mockgen -source=../app/proxy/tablewriter/tablewriterproxy.go -destination=../mock/app/proxy/tablewriter/tablewriterproxy.go -package=mocktablewriterproxy
mockgen -source=../app/proxy/tablewriter/tableinstance.go -destination=../mock/app/proxy/tablewriter/tableinstance.go -package=mocktablewriterproxy
mockgen -source=../app/proxy/time/timeproxy.go -destination=../mock/app/proxy/time/timeproxy.go -package=mocktimeproxy
mockgen -source=../app/proxy/time/timeinstance.go -destination=../mock/app/proxy/time/timeinstance.go -package=mocktimeproxy
mockgen -source=../app/proxy/user/userproxy.go -destination=../mock/app/proxy/user/userproxy.go -package=mockuserproxy
mockgen -source=../app/proxy/user/userinstance.go -destination=../mock/app/proxy/user/userinstance.go -package=mockuserproxy
