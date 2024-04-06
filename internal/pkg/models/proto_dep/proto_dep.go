package proto_dep

type ProtoDep struct {
  Source  string           `json:"source"`
  Targets []ProtoDepTarget `json:"targets"`
}

type ProtoDepTarget struct {
  Source string `json:"source"`
  Import string `json:"import"`
  Commit string `json:"commit"`
}
