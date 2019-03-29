# arrsup
- arrest suppressor

## description
- Suppress infinite redirects when
    - method == "GET" &&
    - origin != host &&
    - set-cookie != null &&
    - status_code / 100 == 3

## motivation
- Want to suppress inifinite redirects while signing in to NextCloud on iOS with SAML.

## License
- Under the MIT License
- Copyright (c) 2019 Tsuzu