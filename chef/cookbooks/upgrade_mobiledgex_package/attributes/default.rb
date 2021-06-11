# MobiledgeX package version
defaultMobiledgeXPackageVersion = '4.4.2'

if node.normal['tags'].include?('vmtype/rootlb')
  if !node.attribute?("mobiledgeXPackageVersion") || node.normal["mobiledgeXPackageVersion"] != defaultMobiledgeXPackageVersion
    node.normal["mobiledgeXPackageVersion"] = defaultMobiledgeXPackageVersion
  end
end
