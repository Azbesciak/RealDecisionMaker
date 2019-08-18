# RealDecisionMaker
Simulator of real human-like decision maker, focused on pairwise comparisons


## Resources
- [Decision Making: Factors that Influence Decision Making, Heuristics Used, and Decision Outcomes](http://www.inquiriesjournal.com/articles/180/decision-making-factors-that-influence-decision-making-heuristics-used-and-decision-outcomes)
- [Naturalistic Decision Making](http://www.decisionskills.com/uploads/5/1/6/0/5160560/klein_2008.pdf)
- [Factors Influencing Familial Decision-Making Regarding Human   Papillomavirus Vaccination](https://watermark.silverchair.com/jsp108.pdf?token=AQECAHi208BE49Ooan9kkhW_Ercy7Dm3ZL_9Cf3qfKAc485ysgAAAkIwggI-BgkqhkiG9w0BBwagggIvMIICKwIBADCCAiQGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMctPCcJhrXoNGq0RdAgEQgIIB9bskWvpXE6eXYEVSejfciigqcuRMWcN_E3lJyQNApzgAPJFX0z4XD2H8gzJeloZ13WT2ph0OUo8iekRkptI0kvMK89P6_oH67DscMXVsl6XWTXgKp-ZbnImsfnVJ9m3o40eMrCTJbD3Fx0Yiq3IlvPZB2Omzv-smR917Gqr9ANOL-77NVedu_RGT165xuEA9Wtv-V8cyMuNgRyiy8_0HG_TBRF3icMqh-dSMp8hfsjXhi5n2V-j7bsgsbh3LAv1kDU-9Y2EJcPlerG45-GEsUBLCNqLGLsRiUVFKf86dJg-mY-g6Vgt_BDMYjx5aHw21f13qJEBsFlHSrj4TxMJNSz8LLOhOVqBUxFy-XhZTorwShRqlOMNrU4EATmBq5vjGl6_yQ5GunG7ulxVbuE8aDBVq_NDpzhHps2WAi6BwX2IWUT-CHwNTmuwelfWQn2BsTB49MNIlX3lmCouK0vAsts-sDhJD5kn3VfHJRYKo7OomcOZtrEtkgEaV10O-zp4mPKsu_iy2pN8q1oOnPTo67_VVaxBU_EbhhcZ_XDWsg28JZL3PbOGOKT95clBgNRhIyOD0RMfyR3gC_Zvw_yu2n8UlOkYouR7Fxv8b4GUZIpXd0jNqcCPlfFIKXO2Zqaq2TnS4FPGP6zLMo-1USWc3-ULXGN93sg)
- [Human Factors Influencing Decision Making](https://apps.dtic.mil/dtic/tr/fulltext/u2/a351910.pdf)
- [Modeling Human and Organizational Behavior: Application to Military Simulations (1998)](https://www.nap.edu/read/6173/chapter/8)
- [Toward a Synthesis of Cognitive Biases: How Noisy Information Processing Can Bias Human Decision Making](https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.432.8763&rep=rep1&type=pdf)

## Thoughts
- DM needs to have some weighting function;
  - Chebyshev
  - Weighted sum
  - Choquet's integral
- what with incomparability? functions above can't manage it (Electre? Promethee I?)
- how to manage "I don't want to answer" / "I don't know"
- Need to remember - decision depends on presented alternatives now (also multiple pairs) and earlier (order is also important).
