import jenkins.model.*
import hudson.security.*
import com.cloudbees.plugins.credentials.*
import org.yaml.snakeyaml.*

DumperOptions options = new DumperOptions()
options.setDefaultFlowStyle(DumperOptions.FlowStyle.BLOCK)

Yaml yaml = new Yaml(options)

def instance = Jenkins.getInstance()
def strategy = new hudson.security.ProjectMatrixAuthorizationStrategy()

class BuildPermission {
    static buildNewAccessList(userOrGroup, permissions) {
        def newPermissionsMap = [:]

        permissions.each {
            println "Permission $it u $userOrGroup perm" + Permission.fromId(it)  
            newPermissionsMap.put(Permission.fromId(it), userOrGroup)
        }

        newPermissionsMap
    }   
}

authenticatedPermissions = [
    "com.cloudbees.plugins.credentials.CredentialsProvider.View",
    "hudson.model.Hudson.Read",
    "hudson.model.Item.Build",
    "hudson.model.Item.Cancel",
    "hudson.model.Item.Discover",
    "hudson.model.Item.Read",
    "hudson.model.Item.Release",
    "hudson.model.Item.Workspace",
    "hudson.model.Run.Delete",
    "hudson.model.Run.Update",
    "hudson.model.View.Configure",
    "hudson.model.View.Create",
    "hudson.model.View.Delete",
    "hudson.model.View.Read"
]

administratorPermissions = [
    "hudson.model.Hudson.Administer"
]

def jenkinsPermissions = yaml.load(("/var/vcap/store/jenkins-master/init.groovy.d/permissions.yml" as File).text)

jenkinsPermissions.each { permissionType ->

    def permissions
    switch (permissionType.name) {
        case "administrator":
          permissions = administratorPermissions
          break
        case "authenticated":
          permissions = authenticatedPermissions
          break
        default:
          println "Unknown permissions type ${permissionType.name}"
          System.exit(1)
    }

    permissionType.organisations.each { organisation ->
       organisation.groups.each { group ->
          BuildPermission.buildNewAccessList("${organisation.name}*${group}", permissions).each { p, u -> strategy.add(p, u) }
        }

        organisation.users.each { user ->
          BuildPermission.buildNewAccessList("${user}", permissions).each { p, u -> strategy.add(p, u) }
        }
    }
}

instance.authorizationStrategy = strategy
instance.save()
