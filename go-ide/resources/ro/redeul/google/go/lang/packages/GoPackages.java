package ro.redeul.google.go.lang.packages;

import com.intellij.openapi.components.AbstractProjectComponent;
import com.intellij.openapi.components.ServiceManager;
import com.intellij.openapi.project.Project;
import com.intellij.openapi.projectRoots.Sdk;
import com.intellij.openapi.roots.OrderRootType;
import com.intellij.openapi.roots.ex.ProjectRootManagerEx;
import com.intellij.openapi.util.NotNullLazyKey;
import com.intellij.openapi.vfs.VirtualFile;
import com.intellij.psi.PsiElement;
import com.intellij.psi.PsiManager;
import com.intellij.reference.SoftReference;
import com.intellij.util.ConcurrencyUtil;
import com.intellij.util.containers.ConcurrentHashMap;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;
import ro.redeul.google.go.lang.psi.GoFile;
import ro.redeul.google.go.lang.psi.GoPackage;
import ro.redeul.google.go.lang.psi.impl.GoPackageImpl;
import ro.redeul.google.go.lang.psi.types.struct.GoTypeStructAnonymousField;

import java.util.concurrent.ConcurrentMap;

import static ro.redeul.google.go.lang.psi.utils.GoPsiUtils.getAs;

public class GoPackages extends AbstractProjectComponent {

    private static final NotNullLazyKey<GoPackages, Project> INSTANCE_KEY = ServiceManager.createLazyKey(GoPackages.class);

    private volatile SoftReference<ConcurrentMap<String, GoPackage>> packagesCache;

    public GoPackages(Project project) {
        super(project);
    }

    @NotNull
    @Override
    public String getComponentName() {
        return "GoPackages";
    }

    public static GoPackages getInstance(Project project) {
        return INSTANCE_KEY.getValue(project);
    }

    /**
     * @param path a fully qualified path name (not a relative one) accessible from one of the source roots.
     * @return an object encapsulating a go package.
     */
    @Nullable
    public GoPackage getPackage(String path) {
        ConcurrentMap<String, GoPackage> cache = SoftReference.dereference(packagesCache);
        if (cache == null) {
            packagesCache = new SoftReference<ConcurrentMap<String, GoPackage>>(cache = new ConcurrentHashMap<String, GoPackage>());
        }

        GoPackage aPackage = cache.get(path);
        if (aPackage != null) {
            return aPackage;
        }

        aPackage = resolvePackage(path);
        return aPackage != null ? ConcurrencyUtil.cacheOrGet(cache, path, aPackage) : null;
    }

    private GoPackage resolvePackage(String path){
        VirtualFile sourceRoots[] = ProjectRootManagerEx.getInstanceEx(myProject).getContentSourceRoots();

        for (VirtualFile sourceRoot : sourceRoots) {
            VirtualFile packagePath = sourceRoot.findFileByRelativePath(path);
            if ( packagePath != null && packagePath.isDirectory()) {
                return new GoPackageImpl(packagePath, sourceRoot, PsiManager.getInstance(myProject));
            }
        }

        Sdk projectSdk = ProjectRootManagerEx.getInstanceEx(myProject).getProjectSdk();
        if ( projectSdk == null )
            return null;

        VirtualFile[] sdkSourceRoots = projectSdk.getRootProvider().getFiles(OrderRootType.SOURCES);

        for (VirtualFile sourceRoot : sdkSourceRoots) {
            VirtualFile packagePath = sourceRoot.findFileByRelativePath(path);
            if ( packagePath != null && packagePath.isDirectory()) {
                return new GoPackageImpl(packagePath, sourceRoot, PsiManager.getInstance(myProject));
            }
        }

        return null;
    }

    public GoPackage getBuiltinPackage() {
        return getPackage("builtin");
    }

    @Nullable
    public static GoPackage getPackageFor(@Nullable PsiElement element) {
        if ( element == null )
            return null;

        GoPackages goPackages = getInstance(element.getProject());

        GoFile goFile = getAs(GoFile.class, element.getContainingFile());

        if ( goFile == null )
            return null;

        String packageImportPath = goFile.getPackageImportPath();

        return goPackages.getPackage(packageImportPath);
    }
}
